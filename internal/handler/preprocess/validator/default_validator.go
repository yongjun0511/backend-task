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

func mustValidateUserFieldDefinitions() map[domain.FieldType]*regexp.Regexp {
	pats := make(map[domain.FieldType]*regexp.Regexp, len(domain.UserFieldDefinitions))

	for ft, meta := range domain.UserFieldDefinitions {
		if meta.Start < 0 || meta.End < 0 || meta.Start >= meta.End {
			panic(fmt.Sprintf("field range for %s: start=%d, end=%d", ft, meta.Start, meta.End))
		}

		if strings.TrimSpace(meta.RegexExpr) == "" {
			panic(fmt.Sprintf("missing regex expression for field %s", ft))
		}

		re := regexp.MustCompile(meta.RegexExpr)
		pats[ft] = re
	}

	return pats
}

func NewDefaultValidator() *DefaultValidator {
	pats := mustValidateUserFieldDefinitions()
	return &DefaultValidator{patterns: pats}
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
