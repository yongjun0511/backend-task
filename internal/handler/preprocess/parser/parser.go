package parser

import "banksalad-backend-task/internal/domain"

type Parser interface {
	ParseLine(line string) (*domain.UserChannelDTO, error)
}
