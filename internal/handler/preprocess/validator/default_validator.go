package validator

import (
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

func NewDefaultValidator() *DefaultValidator {
	pats := make(map[domain.FieldType]*regexp.Regexp, len(domain.UserFieldDefinitions))
	for ft, meta := range domain.UserFieldDefinitions {
		pats[ft] = regexp.MustCompile(meta.RegexExpr)
	}
	return &DefaultValidator{patterns: pats}
}

func (v *DefaultValidator) ValidateLine(line string) (bool, error) {
	for ft, meta := range domain.UserFieldDefinitions {

		if meta.End > len(line) {
			return false, errors.Wrapf(ErrMalformedDataFormat, "data :  %s", line)
		}

		raw := strings.TrimSpace(line[meta.Start:meta.End])
		if !v.patterns[ft].MatchString(raw) {
			return false, errors.Wrapf(ErrInvalidFieldConstraint, "field :  %s", ft)
		}
	}
	return true, nil
}
