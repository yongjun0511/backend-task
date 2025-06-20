package parser

import (
	"strings"

	"banksalad-backend-task/internal/domain"
)

type DefaultParser struct{}

func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

func (p *DefaultParser) ParseLine(line string) ([]domain.ChannelDTO, bool) {

	emailMeta := domain.UserFieldDefinitions[domain.EmailField]
	phoneMeta := domain.UserFieldDefinitions[domain.PhoneField]
	scoreMeta := domain.UserFieldDefinitions[domain.ScoreUpField]

	lineLen := len(line)
	if scoreMeta.End > lineLen {
		return nil, false
	}

	score := strings.TrimSpace(line[scoreMeta.Start:scoreMeta.End])
	if score != "Y" {
		return nil, false
	}

	var result []domain.ChannelDTO

	if emailMeta.End <= lineLen {
		email := strings.TrimSpace(line[emailMeta.Start:emailMeta.End])
		result = append(result, domain.ChannelDTO{
			FieldType: domain.EmailField,
			Value:     email,
		})
	}

	if phoneMeta.End <= lineLen {
		phone := strings.TrimSpace(line[phoneMeta.Start:phoneMeta.End])
		result = append(result, domain.ChannelDTO{
			FieldType: domain.PhoneField,
			Value:     phone,
		})
	}

	return result, true
}
