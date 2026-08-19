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

	"go-reco/projects"
)

func TestFetchProjectsOnce(t *testing.T) {
	t.Setenv("PAGE_SIZE", "1")
	t.Setenv("WORKSPACE_GID", "workspace-gid")

	pages := [][]projects.Project{
		{{Gid: "1", Name: "Roadmap", ResourceType: "project"}},
		{{Gid: "2", Name: "Backlog", ResourceType: "project"}},
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

	projects.Dir = t.TempDir()
	client := resty.New().SetBaseURL(server.URL)

	if err := fetchProjectsOnce(client, time.Second); err != nil {
		t.Fatalf("fetchProjectsOnce() error = %v", err)
	}

	if requests != len(pages) {
		t.Errorf("requests = %d, want %d", requests, len(pages))
	}

	for _, page := range pages {
		for _, project := range page {
			data, err := os.ReadFile(filepath.Join(projects.Dir, project.Gid+".json"))
			if err != nil {
				t.Fatalf("ReadFile(%s) error = %v", project.Gid, err)
			}

			var got projects.Project
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal error = %v", err)
			}
			if got != project {
				t.Errorf("saved project = %+v, want %+v", got, project)
			}
		}
	}
}

func TestFetchProjectsOnceAPIError(t *testing.T) {
	t.Setenv("PAGE_SIZE", "1")
	t.Setenv("WORKSPACE_GID", "workspace-gid")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	projects.Dir = t.TempDir()
	client := resty.New().SetBaseURL(server.URL)

	if err := fetchProjectsOnce(client, time.Second); err == nil {
		t.Fatal("fetchProjectsOnce() error = nil, want error")
	}
}
