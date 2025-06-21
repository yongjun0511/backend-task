package preprocess

import (
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"
	"bufio"
	"github.com/pkg/errors"
	"os"

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

func (pp *Preprocessor) Run() (map[domain.FieldType]map[string]struct{}, error) {
	f, err := os.Open(pp.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[domain.FieldType]map[string]struct{})

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		if err := pp.validator.ValidateLine(line); err != nil {
			if errors.Is(err, validator.ErrMalformedDataFormat) ||
				errors.Is(err, validator.ErrInvalidFieldConstraint) {
				logrus.WithError(err).Warn("skip record during validation")
				continue
			}
			return nil, errors.WithStack(err)
		}

		dto, err := pp.parser.ParseLine(line)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if dto == nil {
			continue
		}

		if dto.Email != "" {
			if _, ok := result[domain.EmailField]; !ok {
				result[domain.EmailField] = make(map[string]struct{})
			}
			result[domain.EmailField][dto.Email] = struct{}{}
		}

		if dto.SMS != "" {
			if _, ok := result[domain.PhoneField]; !ok {
				result[domain.PhoneField] = make(map[string]struct{})
			}
			result[domain.PhoneField][dto.SMS] = struct{}{}
		}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
