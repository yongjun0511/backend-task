package parser

import "banksalad-backend-task/internal/domain"

type Parser interface {
	ParseLine(path string) ([]domain.UserRecord, error)
}
