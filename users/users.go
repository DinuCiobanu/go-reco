package users

import (
	"encoding/json"
	"os"
	"path/filepath"

	"resty.dev/v3"

	"go-reco/api"
)

const Dir = "data/users"

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

func FetchAll(client *resty.Client, onPage func([]User)) error {
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

		resp, err := api.GetWithRetry(req, "/users")
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

func Save(user User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	path := filepath.Join(Dir, user.Gid+".json")
	return os.WriteFile(path, data, 0o644)
}
