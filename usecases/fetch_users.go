package usecases

import (
	"fmt"
	"os"
	"time"

	"resty.dev/v3"

	"go-reco/api"
	"go-reco/users"
	"go-reco/utilities"
)

func FetchUsersAtEvery(duration int, unit time.Duration) {
	client := api.AsanaClient()
	interval := time.Duration(duration) * unit

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := fetchUsersOnce(client, interval); err != nil {
			fmt.Println("error:", err)
		}

		<-ticker.C
	}
}

func fetchUsersOnce(client *resty.Client, interval time.Duration) error {
	if err := os.MkdirAll(users.Dir, 0o755); err != nil {
		return err
	}

	utilities.LogFetch("users", interval)

	return users.FetchAll(client, func(fetched []users.User) {
		for _, user := range fetched {
			if err := users.Save(user); err != nil {
				fmt.Println("error saving user:", err)
			}
		}
	})
}
