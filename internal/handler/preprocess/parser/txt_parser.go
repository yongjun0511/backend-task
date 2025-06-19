package parser

import (
	"strings"

	"banksalad-backend-task/internal/domain"
)

type TxtParser struct{}

func (p *TxtParser) ParseLine(line string) ([]domain.ChannelDTO, bool) {

	fields := strings.Fields(line)

	if len(fields) != 3 {
		return nil, false
	}

	email := fields[0]
	phone := fields[1]
	scoreUp := fields[2] == "Y"

	if !scoreUp {
		return nil, false
	}

	return []domain.ChannelDTO{
		{Channel: domain.EmailChannel, Value: email},
		{Channel: domain.SMSChannel, Value: phone},
	}, true
}
