package channelhandler

import (
	"log"
	"sync"

	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"
)

type EmailHandler struct {
	client *clients.EmailClient
	mu     sync.Mutex
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		client: clients.NewEmailClient(),
	}
}

func (h *EmailHandler) TargetField() domain.FieldType {
	return domain.EmailField
}

func (h *EmailHandler) SendBatch(values []string) error {
	for _, email := range values {
		for {
			h.mu.Lock()
			err := h.client.Send(email, "신용 점수가 상승했습니다!")
			h.mu.Unlock()

			if err == nil {
				break
			}
			log.Printf("[WARN] 이메일 전송 실패: %s, 재시도 중... err: %v", email, err)
		}
	}
	return nil
}
