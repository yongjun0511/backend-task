package channelhandler

import (
	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"
	"log"
	"sync"
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

func (h *EmailHandler) SendBatch(values []string) error {
	var wg sync.WaitGroup
	for _, v := range values {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			for {
				if err := h.client.Send(email, "신용 점수가 상승했습니다!"); err == nil {
					return
				}
				log.Printf("[WARN] 이메일 전송 실패. 재시도 중... 대상: %s", email)
			}
		}(v)
	}
	wg.Wait()
	return nil
}
