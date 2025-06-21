package parser

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"banksalad-backend-task/internal/domain"
)

var ErrUnexpectedParsing = errors.New("unsafe parsing range, validation bypassed")

type DefaultParser struct{}

func mustValidateFieldRanges() {
	for _, ft := range []domain.FieldType{
		domain.EmailField,
		domain.PhoneField,
		domain.ScoreUpField,
	} {
		meta := domain.UserFieldDefinitions[ft]
		if meta.Start < 0 || meta.End < 0 || meta.Start >= meta.End {
			panic(fmt.Sprintf("field range for %s: start=%d, end=%d",
				ft, meta.Start, meta.End))
		}
	}
}

func NewDefaultParser() *DefaultParser {
	mustValidateFieldRanges()
	return &DefaultParser{}
}
func (p *DefaultParser) ParseLine(line string) (*domain.UserChannelDTO, error) {
	emailMeta := domain.UserFieldDefinitions[domain.EmailField]
	phoneMeta := domain.UserFieldDefinitions[domain.PhoneField]
	scoreMeta := domain.UserFieldDefinitions[domain.ScoreUpField]

	lineLen := len(line)
	for ft, meta := range domain.UserFieldDefinitions {
		if meta.End > lineLen || meta.Start > meta.End {
			return nil, errors.Wrapf(ErrUnexpectedParsing,
				"range invalid: field=%s lineLen=%d meta=%+v line=%q",
				ft, lineLen, meta, line)
		}
	}

	score := strings.TrimSpace(line[scoreMeta.Start:scoreMeta.End])
	if score != "Y" {
		return &domain.UserChannelDTO{}, nil
	}

	email := strings.TrimSpace(line[emailMeta.Start:emailMeta.End])
	phone := strings.TrimSpace(line[phoneMeta.Start:phoneMeta.End])

	return &domain.UserChannelDTO{
		Email: email,
		SMS:   phone,
	}, nil
}
