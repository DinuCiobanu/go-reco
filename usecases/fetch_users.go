package usecases

import (
	"fmt"
	"os"
	"time"

	"go-reco/api"
	"go-reco/users"
	"go-reco/utilities"
)

func FetchUsersAtEvery(duration int, unit time.Duration) {
	if err := os.MkdirAll(users.Dir, 0o755); err != nil {
		fmt.Println("error:", err)
		return
	}

	client := api.AsanaClient()

	onPage := func(fetched []users.User) {
		for _, user := range fetched {
			if err := users.Save(user); err != nil {
				fmt.Println("error saving user:", err)
			}
		}
	}

	interval := time.Duration(duration) * unit
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		utilities.LogFetch("users", interval)
		if err := users.FetchAll(client, onPage); err != nil {
			fmt.Println("error:", err)
		}

		<-ticker.C
	}
}
