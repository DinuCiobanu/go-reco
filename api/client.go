package api

import (
	"os"
	"sync"

	"resty.dev/v3"
)

var (
	client     *resty.Client
	clientOnce sync.Once
)

func AsanaClient() *resty.Client {
	clientOnce.Do(func() {
		client = resty.New().
			SetBaseURL("https://app.asana.com/api/1.0").
			SetAuthToken(os.Getenv("ACCESS_TOKEN"))
	})
	return client
}
