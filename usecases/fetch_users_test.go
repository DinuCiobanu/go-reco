package usecases

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"resty.dev/v3"

	"go-reco/users"
)

func TestFetchUsersOnce(t *testing.T) {
	t.Setenv("PAGE_SIZE", "1")
	t.Setenv("WORKSPACE_GID", "workspace-gid")

	pages := [][]users.User{
		{{Gid: "1", Name: "Alice", ResourceType: "user"}},
		{{Gid: "2", Name: "Bob", ResourceType: "user"}},
	}

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("workspace"); got != "workspace-gid" {
			t.Errorf("workspace query param = %q, want %q", got, "workspace-gid")
		}

		page := pages[requests]
		requests++

		resp := map[string]any{"data": page}
		if requests < len(pages) {
			resp["next_page"] = map[string]string{"offset": "next"}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	users.Dir = t.TempDir()
	client := resty.New().SetBaseURL(server.URL)

	if err := fetchUsersOnce(client, time.Second); err != nil {
		t.Fatalf("fetchUsersOnce() error = %v", err)
	}

	if requests != len(pages) {
		t.Errorf("requests = %d, want %d", requests, len(pages))
	}

	for _, page := range pages {
		for _, user := range page {
			data, err := os.ReadFile(filepath.Join(users.Dir, user.Gid+".json"))
			if err != nil {
				t.Fatalf("ReadFile(%s) error = %v", user.Gid, err)
			}

			var got users.User
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal error = %v", err)
			}
			if got != user {
				t.Errorf("saved user = %+v, want %+v", got, user)
			}
		}
	}
}

func TestFetchUsersOnceAPIError(t *testing.T) {
	t.Setenv("PAGE_SIZE", "1")
	t.Setenv("WORKSPACE_GID", "workspace-gid")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	users.Dir = t.TempDir()
	client := resty.New().SetBaseURL(server.URL)

	if err := fetchUsersOnce(client, time.Second); err == nil {
		t.Fatal("fetchUsersOnce() error = nil, want error")
	}
}
