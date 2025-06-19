package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type ValidationError string

const (
	ErrEmptyLine     ValidationError = "빈줄입니다"
	ErrInvalidLength ValidationError = "데이터의 형식과 길이 올바르지 않습니다)"
	ErrInvalidEmail  ValidationError = "이메일 형식이 올바르지 않습니다"
	ErrInvalidPhone  ValidationError = "전화번호 형식이 올바르지 않습니다 (000-0000-0000)"
	ErrInvalidScore  ValidationError = "신용 등급 정보 형식이 올바르지 않습니다."
)

var (
	emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	phoneRegex = regexp.MustCompile(`^\d{3}-\d{4}-\d{4}$`)
	scoreRegex = regexp.MustCompile(`^[YN]$`)
)

type TxtValidator struct{}

const (
	lineLength       = 72
	emailLengthLimit = 50
	phoneNumberLimit = 72
	scoreStart       = 72
)

func (v *TxtValidator) ValidateLine(line string) (bool, ValidationError) {
	trimRight := strings.TrimRight(line, " \r\n")

	if len(strings.TrimSpace(trimRight)) == 0 {
		return false, ErrEmptyLine
	}

	if utf8.RuneCountInString(trimRight) != lineLength {
		return false, ErrInvalidLength
	}

	email := strings.TrimSpace(trimRight[:emailLengthLimit])
	phone := strings.TrimSpace(trimRight[emailLengthLimit:phoneNumberLimit])
	score := strings.TrimSpace(trimRight[scoreStart:])

	if !emailRegex.MatchString(email) {
		return false, ErrInvalidEmail
	}

	if !phoneRegex.MatchString(phone) {
		return false, ErrInvalidPhone
	}

	if !scoreRegex.MatchString(score) {
		return false, ErrInvalidScore
	}

	return true, ""
}
