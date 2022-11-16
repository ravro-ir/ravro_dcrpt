package utils

import (
	"github.com/google/logger"
	"log"
	"os"
)

func Logger(msg error) {
	lf, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	if err != nil {
		log.Fatal(err)
	}

	defer logger.Init("Error logger : ", true, false, lf).Close()
	logger.SetFlags(log.Llongfile)
	logger.Errorf("Error running: %v", msg.Error())

}
