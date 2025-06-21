package channelhandler

import (
	"sync"
	"time"

	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maxSmsRetry     = 3
	rateLimitPerSec = 100
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

func (h *SMSHandler) TargetField() domain.FieldType {
	return domain.PhoneField
}

func (h *SMSHandler) SendBatch(values []string) error {
	ticker := time.NewTicker(time.Second / rateLimitPerSec)
	defer ticker.Stop()

	for _, phone := range values {
		<-ticker.C

		var err error
		for i := 0; i < maxSmsRetry; i++ {
			h.mu.Lock()
			err = h.client.Send(phone, "신용 점수가 상승했습니다!")
			h.mu.Unlock()

			if err == nil {
				break
			}

			logrus.WithFields(logrus.Fields{
				"attempt": i + 1,
				"max":     maxSmsRetry,
				"phone":   phone,
				"error":   err,
			}).Warn("sms send failed")

			time.Sleep(10 * time.Millisecond)
		}

		if err != nil {
			return errors.WithStack(errors.Wrapf(err, "sms send failed after %d attempts (%s)", maxSmsRetry, phone))
		}
	}

	return nil
}
