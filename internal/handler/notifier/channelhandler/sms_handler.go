package channelhandler

import (
	"log"
	"sync"
	"time"

	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"
)

const (
	tokenPerSec = 100
)

type SMSHandler struct {
	client *clients.SmsClient
	mu     sync.Mutex
}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{
		client: clients.NewSmsClient(),
	}
}

func (h *SMSHandler) TargetField() domain.FieldType { return domain.PhoneField }

func (h *SMSHandler) SendBatch(values []string) error {
	queue := make(chan string, len(values)*2)
	for _, v := range values {
		queue <- v
	}

	tokenCh := make(chan struct{}, tokenPerSec)
	go refillTokens(tokenCh)

	for i := 0; i < tokenPerSec; i++ {
		go h.worker(queue, tokenCh)
	}

	for len(queue) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func refillTokens(tokenCh chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		for i := 0; i < tokenPerSec; i++ {
			tokenCh <- struct{}{}
		}
	}
}

func (h *SMSHandler) worker(queue chan string, tokenCh chan struct{}) {
	for {
		select {
		case phone := <-queue:
			<-tokenCh
			for {
				h.mu.Lock()
				err := h.client.Send(phone, "신용 점수가 상승했습니다!")
				h.mu.Unlock()

				if err == nil {
					break
				}
				log.Printf("[WARN] SMS 실패 → 다음 초 재시도: %s", phone)
				time.Sleep(10 * time.Millisecond)
			}
		case <-time.After(2 * time.Second):
			if len(queue) == 0 {
				return
			}
		}
	}
}
