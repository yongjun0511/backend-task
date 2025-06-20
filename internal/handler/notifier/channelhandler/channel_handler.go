package channelhandler

import "banksalad-backend-task/internal/domain"

type ChannelHandler interface {
	TargetField() domain.FieldType
	SendBatch(values []string) error
}
