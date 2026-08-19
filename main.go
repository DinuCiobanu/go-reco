package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"resty.dev/v3"
)

const usersDir = "data/users"

func newAsanaClient() *resty.Client {
	return resty.New().
		SetBaseURL("https://app.asana.com/api/1.0").
		SetAuthToken(os.Getenv("ACCESS_TOKEN"))
}

type User struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

type usersPage struct {
	Data     []User `json:"data"`
	NextPage *struct {
		Offset string `json:"offset"`
	} `json:"next_page"`
}

func getWithRetry(req *resty.Request, url string) (*resty.Response, error) {
	for {
		resp, err := req.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() != http.StatusTooManyRequests {
			return resp, nil
		}

		wait := 1
		if seconds, err := strconv.Atoi(resp.Header().Get("Retry-After")); err == nil {
			wait = seconds
		}

		fmt.Printf("rate limited, retrying in %ds\n", wait)
		time.Sleep(time.Duration(wait) * time.Second)
	}
}

func fetchAllUsers(client *resty.Client, onPage func([]User)) error {
	pageSize := os.Getenv("PAGE_SIZE")

	workspace := os.Getenv("WORKSPACE_GID")

	offset := ""

	for {
		req := client.R().
			SetQueryParam("workspace", workspace).
			SetQueryParam("limit", pageSize)
		if offset != "" {
			req.SetQueryParam("offset", offset)
		}

		resp, err := getWithRetry(req, "/users")
		if err != nil {
			return err
		}

		var page usersPage
		if err := json.Unmarshal(resp.Bytes(), &page); err != nil {
			return err
		}

		onPage(page.Data)

		if page.NextPage == nil || page.NextPage.Offset == "" {
			break
		}
		offset = page.NextPage.Offset
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found, falling back to environment")
	}

	if err := os.MkdirAll(usersDir, 0o755); err != nil {
		fmt.Println("error:", err)
		return
	}

	client := newAsanaClient()

	err := fetchAllUsers(client, func(users []User) {
		for _, user := range users {
			if err := saveUser(user); err != nil {
				fmt.Println("error saving user:", err)
			}
		}
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

func saveUser(user User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	path := filepath.Join(usersDir, user.Gid+".json")
	return os.WriteFile(path, data, 0o644)
}
