package preprocess

import (
	"bufio"
	"os"

	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/parser"
	"banksalad-backend-task/internal/validator"
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

func (pp *Preprocessor) Run() (map[domain.ChannelType]map[string]struct{}, error) {
	f, err := os.Open(pp.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 🔑 키 타입을 domain.ChannelType으로!
	result := make(map[domain.ChannelType]map[string]struct{})

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		vals, ok := pp.parser.ParseLine(line)
		if !ok {
			continue
		}

		for _, cv := range vals {
			if _, exists := result[cv.Channel]; !exists {
				result[cv.Channel] = make(map[string]struct{})
			}
			result[cv.Channel][cv.Value] = struct{}{}
		}
	}
	return result, sc.Err()
}
