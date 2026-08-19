package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go-reco/usecases"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found, falling back to environment")
	}

	duration, unit := resolveInterval()

	go usecases.FetchUsersAtEvery(duration, unit)
	usecases.FetchProjectsAtEvery(duration, unit)
}

func resolveInterval() (int, time.Duration) {
	switch os.Getenv("FETCHING_MODE") {
	case "2":
		return 5, time.Minute
	default:
		return 30, time.Second
	}
}
