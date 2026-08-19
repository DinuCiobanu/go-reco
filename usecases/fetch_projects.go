package usecases

import (
	"fmt"
	"os"
	"time"

	"resty.dev/v3"

	"go-reco/api"
	"go-reco/projects"
	"go-reco/utilities"
)

func FetchProjectsAtEvery(duration int, unit time.Duration) {
	client := api.AsanaClient()
	interval := time.Duration(duration) * unit

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := fetchProjectsOnce(client, interval); err != nil {
			fmt.Println("error:", err)
		}

		<-ticker.C
	}
}

func fetchProjectsOnce(client *resty.Client, interval time.Duration) error {
	if err := os.MkdirAll(projects.Dir, 0o755); err != nil {
		return err
	}

	utilities.LogFetch("projects", interval)

	return projects.FetchAll(client, func(fetched []projects.Project) {
		for _, project := range fetched {
			if err := projects.Save(project); err != nil {
				fmt.Println("error saving project:", err)
			}
		}
	})
}
