package filter

import "banksalad-backend-task/internal/domain"

type DefaultContactFilter struct{}

func (f *DefaultContactFilter) Extract(records []domain.UserRecord) (map[string]struct{}, map[string]struct{}) {
	emailSet := make(map[string]struct{})
	phoneSet := make(map[string]struct{})

	for _, r := range records {
		if !r.ScoreUp {
			continue
		}
		emailSet[r.Email] = struct{}{}
		phoneSet[r.Phone] = struct{}{}
	}

	return emailSet, phoneSet
}
