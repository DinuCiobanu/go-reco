package api

import (
	"os"

	"resty.dev/v3"
)

func AsanaClient() *resty.Client {
	return resty.New().
		SetBaseURL("https://app.asana.com/api/1.0").
		SetAuthToken(os.Getenv("ACCESS_TOKEN"))
}
