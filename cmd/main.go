package main

import (
	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/handler/preprocess"
	"log"
	"os"

	"banksalad-backend-task/internal/handler/notifier"
	"banksalad-backend-task/internal/handler/notifier/channelhandler"
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"

	"github.com/sirupsen/logrus"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("reason", r).Error("💥 application initialization failed")
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

	result, err := pp.Run()
	if err != nil {
		logrus.WithError(err).Fatal("preprocessing failed")
	}

	if err := nt.NotifyAll(result); err != nil {
		log.Fatalf("[ERROR] Notification failed: %v", err)
	}

	log.Println("[INFO] All notifications sent successfully.")
}
