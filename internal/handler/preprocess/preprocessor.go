package preprocess

import (
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"
	"bufio"
	"context"
	"github.com/pkg/errors"
	"os"
	"sync"

	"banksalad-backend-task/internal/domain"

	"github.com/sirupsen/logrus"
)

type Preprocessor struct {
	path      string
	parser    parser.Parser
	validator validator.Validator
}

func NewPreprocessor(
	path string,
	p parser.Parser,
	v validator.Validator,
) *Preprocessor {
	return &Preprocessor{
		path:      path,
		parser:    p,
		validator: v,
	}
}
func (pp *Preprocessor) Run(ctx context.Context, workers int) (map[domain.FieldType]map[string]struct{}, error) {
	f, err := os.Open(pp.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	linesCh := make(chan string, workers*4)
	errCh := make(chan error, 1)

	var (
		emailSet sync.Map
		phoneSet sync.Map
		wg       sync.WaitGroup
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range linesCh {
				if ctx.Err() != nil {
					return
				}

				if err := pp.validator.ValidateLine(line); err != nil {
					if errors.Is(err, validator.ErrMalformedDataFormat) ||
						errors.Is(err, validator.ErrInvalidFieldConstraint) {
						logrus.WithError(err).Warn("skip record during validation")
						continue
					}
					select {
					case errCh <- errors.WithStack(err):
					default:
					}
					return
				}

				dto, err := pp.parser.ParseLine(line)
				if err != nil {
					select {
					case errCh <- errors.WithStack(err):
					default:
					}
					return
				}
				if dto == nil {
					continue
				}

				if dto.Email != "" {
					emailSet.Store(dto.Email, struct{}{})
				}
				if dto.SMS != "" {
					phoneSet.Store(dto.SMS, struct{}{})
				}
			}
		}()
	}

	go func() {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			select {
			case <-ctx.Done():
				break
			case linesCh <- sc.Text():
			}
		}
		close(linesCh)
		if err := sc.Err(); err != nil {
			errCh <- err
		}
	}()

	doneCh := make(chan struct{})
	go func() { wg.Wait(); close(doneCh) }()

	select {
	case err := <-errCh:
		return nil, err
	case <-doneCh:
	}

	result := map[domain.FieldType]map[string]struct{}{
		domain.EmailField: make(map[string]struct{}),
		domain.PhoneField: make(map[string]struct{}),
	}
	emailSet.Range(func(k, _ any) bool { result[domain.EmailField][k.(string)] = struct{}{}; return true })
	phoneSet.Range(func(k, _ any) bool { result[domain.PhoneField][k.(string)] = struct{}{}; return true })

	return result, nil
}
