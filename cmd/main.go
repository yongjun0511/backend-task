package main

import (
	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/handler/preprocess"
	"context"
	"log"
	"os"
	"time"

	"banksalad-backend-task/internal/handler/notifier"
	"banksalad-backend-task/internal/handler/notifier/channelhandler"
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"

	"github.com/sirupsen/logrus"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("reason", r).Error("application initialization failed")
			os.Exit(1)
		}
	}()

	pp := preprocess.NewPreprocessor(
		"files/input/data.txt",
		parser.NewDefaultParser(),
		validator.NewDefaultValidator(),
	)

	emailHandler := channelhandler.NewEmailHandler()
	smsHandler := channelhandler.NewSMSHandler()

	nt := notifier.NewNotifier(map[domain.FieldType]channelhandler.ChannelHandler{
		domain.EmailField: emailHandler,
		domain.PhoneField: smsHandler,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pp.Run(ctx, 4)
	if err != nil {
		logrus.WithError(err).Fatal("preprocessing failed")
	}

	if err := nt.NotifyAll(result); err != nil {
		logrus.WithError(err).Fatal("notification failed")
	}

	log.Println("[INFO] All notifications sent successfully.")
}
