package channelhandler

import (
	"fmt"
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
		var err error
		for i := 0; i < 3; i++ {
			h.mu.Lock()
			err = h.client.Send(email, "신용 점수가 상승했습니다!")
			h.mu.Unlock()

			if err == nil {
				break
			}
			log.Printf("[WARN] 이메일 전송 실패 (%d/3): %s, err: %v", i+1, email, err)
		}

		if err != nil {
			return fmt.Errorf("이메일 전송 실패 (email: %s): %w", email, err)
		}
	}
	return nil
}
