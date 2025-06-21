package channelhandler

import (
	"github.com/pkg/errors"
	"sync"

	"banksalad-backend-task/clients"
	"banksalad-backend-task/internal/domain"

	"github.com/sirupsen/logrus"
)

const maxRetry = 3

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
		for i := 0; i < maxRetry; i++ {
			h.mu.Lock()
			err = h.client.Send(email, "신용 점수가 상승했습니다!")
			h.mu.Unlock()

			if err == nil {
				break
			}
			logrus.WithFields(logrus.Fields{
				"attempt": i + 1,
				"max":     maxRetry,
				"email":   email,
				"error":   err,
			}).Warn("email send failed")
		}

		if err != nil {
			return errors.Wrapf(err, "email send failed after %d attempts (%s)", maxRetry, email)
		}
	}
	return nil
}
