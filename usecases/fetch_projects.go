package usecases

import (
	"fmt"
	"os"
	"time"

	"go-reco/api"
	"go-reco/projects"
	"go-reco/utilities"
)

func FetchProjectsAtEvery(duration int, unit time.Duration) {
	if err := os.MkdirAll(projects.Dir, 0o755); err != nil {
		fmt.Println("error:", err)
		return
	}

	client := api.AsanaClient()

	onPage := func(fetched []projects.Project) {
		for _, project := range fetched {
			if err := projects.Save(project); err != nil {
				fmt.Println("error saving project:", err)
			}
		}
	}

	interval := time.Duration(duration) * unit
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		utilities.LogFetch("projects", interval)
		if err := projects.FetchAll(client, onPage); err != nil {
			fmt.Println("error:", err)
		}

		<-ticker.C
	}
}
