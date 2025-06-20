package channelhandler

import (
	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"
)

type EmailHandler struct {
	client *clients.EmailClient
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		client: clients.NewEmailClient(),
	}
}

func (h *EmailHandler) TargetField() domain.FieldType {
	return domain.EmailField
}

func (h *EmailHandler) Send(value string) error {
	return h.client.Send(value, "신용 점수가 상승했습니다!") // 내용 고민해볼 것.
}
