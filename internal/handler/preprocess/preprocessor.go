package preprocess

import (
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"
	"bufio"
	"os"

	"banksalad-backend-task/internal/domain"
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

func (pp *Preprocessor) Run() (map[domain.ChannelDTO]map[string]struct{}, error) {
	f, err := os.Open(pp.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[domain.ChannelDTO]map[string]struct{})

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		isValid, err := pp.validator.ValidateLine(line)
		if err != nil || !isValid {
			continue
		}

		vals, ok := pp.parser.ParseLine(line)
		if !ok {
			continue
		}

		for _, fv := range vals {
			if _, exists := result[fv]; !exists {
				result[fv] = make(map[string]struct{})
			}
			result[fv][fv.Value] = struct{}{}
		}
	}
	return result, sc.Err()
}
