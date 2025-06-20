package main

import (
	"banksalad-backend-task/internal/handler/preprocess"
	"log"

	"banksalad-backend-task/internal/handler/notifier"
	"banksalad-backend-task/internal/handler/notifier/channelhandler"
	"banksalad-backend-task/internal/handler/preprocess/parser"
	"banksalad-backend-task/internal/handler/preprocess/validator"
)

func main() {
	pp := preprocess.NewPreprocessor(
		"files/input/data.txt",
		parser.NewDefaultParser(),
		validator.NewDefaultValidator(),
	)

	result, err := pp.Run()
	if err != nil {
		log.Fatalf("[ERROR] Preprocessing failed: %v", err)
	}

	emailHandler := channelhandler.NewEmailHandler()
	smsHandler := channelhandler.NewSMSHandler()

	nt := notifier.NewNotifier([]channelhandler.ChannelHandler{
		emailHandler,
		smsHandler,
	})

	if err := nt.NotifyAll(result); err != nil {
		log.Fatalf("[ERROR] Notification failed: %v", err)
	}

	log.Println("[INFO] All notifications sent successfully.")
}
