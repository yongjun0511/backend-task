package channelhandler

import (
	"time"

	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"
)

type SMSHandler struct {
	client *clients.SmsClient
}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{
		client: clients.NewSmsClient(),
	}
}

func (h *SMSHandler) TargetField() domain.FieldType {
	return domain.PhoneField
}

func (h *SMSHandler) Send(value string) error {
	time.Sleep(10 * time.Millisecond) // 초당 100건 제한을 단순하게 반영 -> 나중에 수정해야할듯
	return h.client.Send(value, "신용 점수가 상승했습니다!")
}
