package usecases

import (
	"fmt"
	"os"
	"time"

	"go-reco/api"
	"go-reco/users"
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

	ticker := time.NewTicker(time.Duration(duration) * unit)
	defer ticker.Stop()

	for {
		fmt.Println("Fetch users")
		if err := users.FetchAll(client, onPage); err != nil {
			fmt.Println("error:", err)
		}

		<-ticker.C
	}
}
