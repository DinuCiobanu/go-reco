package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"go-reco/usecases"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found, falling back to environment")
	}

	usecases.FetchUsersAtEvery(5, time.Second)
}
