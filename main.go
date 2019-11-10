package main

import (
	"log"

	"github.com/2hamed/saas/jobq"
	"github.com/2hamed/saas/screenshot"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	capture, err := screenshot.NewCapture()
	if err != nil {
		log.Fatal(err)
	}

	_, err = jobq.NewQManager(capture, nil)
	if err != nil {
		log.Fatal(err)
	}
}
