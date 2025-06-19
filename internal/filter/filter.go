package filter

import "banksalad-backend-task/internal/domain"

type ContactFilter interface {
	Extract(records []domain.UserRecord) (map[string]struct{}, map[string]struct{})
}
