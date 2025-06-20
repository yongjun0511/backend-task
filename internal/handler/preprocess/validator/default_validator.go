package validator

import (
	"regexp"
	"strings"

	"banksalad-backend-task/internal/domain"
)

type ValidationError struct {
	Field  domain.FieldType
	Reason string
	Raw    string
}

func (e ValidationError) Error() string {
	return string(e.Field) + ": " + e.Reason + " (" + e.Raw + ")"
}

type DefaultValidator struct {
	patterns map[domain.FieldType]*regexp.Regexp
}

func NewDefaultValidator() *DefaultValidator {
	compiled := make(map[domain.FieldType]*regexp.Regexp, len(domain.UserFieldDefinitions))
	for ft, meta := range domain.UserFieldDefinitions {
		compiled[ft] = regexp.MustCompile(meta.RegexExpr)
		_ = ft
	}
	return &DefaultValidator{patterns: compiled}
}

func (v *DefaultValidator) ValidateLine(line string) (bool, error) {
	for ft, meta := range domain.UserFieldDefinitions {
		if meta.End > len(line) {
			return false, ValidationError{
				Field:  ft,
				Raw:    line,
				Reason: "데이터 길이 오류 ",
			}
		}

		raw := strings.TrimSpace(line[meta.Start:meta.End])
		if !v.patterns[ft].MatchString(raw) {
			return false, ValidationError{
				Field:  ft,
				Raw:    raw,
				Reason: "데이터 형식 오류",
			}
		}
	}
	return true, nil
}
