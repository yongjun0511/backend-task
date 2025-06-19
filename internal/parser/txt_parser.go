package parser

import (
	"banksalad-backend-task/internal/domain"
	"strings"
)

type TxtParser struct{}

func (p *TxtParser) ParseLine(line string) domain.UserRecord {
	fields := strings.Fields(line)

	scoreUp := fields[2] == "Y"

	return domain.UserRecord{
		Email:   fields[0],
		Phone:   fields[1],
		ScoreUp: scoreUp,
	}
}
