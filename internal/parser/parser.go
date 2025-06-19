package parser

import "banksalad-backend-task/internal/domain"

type Parser interface {
	Parse(path string) ([]domain.UserRecord, error)
}
