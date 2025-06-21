package validator

import (
	"fmt"
	"regexp"
	"strings"

	"banksalad-backend-task/internal/domain"

	"github.com/pkg/errors"
)

var (
	ErrMalformedDataFormat    = errors.New("data format")
	ErrInvalidFieldConstraint = errors.New("field constraint")
)

type DefaultValidator struct {
	patterns map[domain.FieldType]*regexp.Regexp
}

func NewDefaultValidator() (*DefaultValidator, error) {
	pats := make(map[domain.FieldType]*regexp.Regexp, len(domain.UserFieldDefinitions))

	for ft, meta := range domain.UserFieldDefinitions {

		// 상수 index check
		if meta.Start < 0 || meta.End < 0 || meta.Start >= meta.End {
			return nil, fmt.Errorf("invalid field range for %s: start=%d, end=%d", ft, meta.Start, meta.End)
		}

		// 정규식 유무 확인
		if strings.TrimSpace(meta.RegexExpr) == "" {
			return nil, fmt.Errorf("missing regex expression for field %s", ft)
		}

		// 3. 정규식 컴파일 (에러 발생 가능성 있음)
		re, err := regexp.Compile(meta.RegexExpr)
		if err != nil {
			return nil, fmt.Errorf("invalid regex for field %s: %v", ft, err)
		}

		pats[ft] = re
	}

	return &DefaultValidator{patterns: pats}, nil
}

func MustDefaultValidator() *DefaultValidator {
	v, err := NewDefaultValidator()
	if err != nil {
		panic(err)
	}
	return v
}

func (v *DefaultValidator) ValidateLine(line string) error {
	for ft, meta := range domain.UserFieldDefinitions {

		if meta.End > len(line) {
			return errors.Wrapf(ErrMalformedDataFormat, "line = %s", line)
		}

		raw := strings.TrimSpace(line[meta.Start:meta.End])
		if !v.patterns[ft].MatchString(raw) {
			return errors.Wrapf(ErrInvalidFieldConstraint, "field = %s", ft)
		}
	}
	return nil
}
