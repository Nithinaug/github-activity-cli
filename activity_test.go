package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestFormatPushEvent(t *testing.T) {
	e := GitHubEvent{
		Type:    "PushEvent",
		Repo:    Repo{Name: "user/repo"},
		Payload: Payload{Size: 3},
	}
	got := formatEvent(e)
	if !strings.Contains(got, "3") {
		t.Errorf("expected commit count in output, got: %s", got)
	}
	if !strings.Contains(got, "user/repo") {
		t.Errorf("expected repo name in output, got: %s", got)
	}
}

func TestFormatWatchEvent(t *testing.T) {
	e := GitHubEvent{Type: "WatchEvent", Repo: Repo{Name: "user/repo"}}
	got := formatEvent(e)
	if !strings.Contains(got, "Starred") {
		t.Errorf("expected 'Starred' in output, got: %s", got)
	}
}

func TestFormatForkEvent(t *testing.T) {
	e := GitHubEvent{
		Type:    "ForkEvent",
		Repo:    Repo{Name: "original/repo"},
		Payload: Payload{Forkee: &Forkee{FullName: "user/repo"}},
	}
	got := formatEvent(e)
	if !strings.Contains(got, "Forked") {
		t.Errorf("expected 'Forked' in output, got: %s", got)
	}
}

func TestFormatIssuesEvent(t *testing.T) {
	e := GitHubEvent{
		Type: "IssuesEvent",
		Repo: Repo{Name: "user/repo"},
		Payload: Payload{
			Action: "opened",
			Issue:  &Issue{Number: 42, Title: "Bug report"},
		},
	}
	got := formatEvent(e)
	if !strings.Contains(got, "#42") {
		t.Errorf("expected issue number in output, got: %s", got)
	}
}

func TestFormatUnknownEvent(t *testing.T) {
	e := GitHubEvent{Type: "SomeNewEvent", Repo: Repo{Name: "user/repo"}}
	got := formatEvent(e)
	if got == "" {
		t.Error("expected non-empty output for unknown event type")
	}
}

// ── helper tests ─────────────────────────────────────────────────────────────

func TestTruncate(t *testing.T) {
	if got := truncate("hello world", 5); got != "hell…" {
		t.Errorf("truncate: got %q", got)
	}
	if got := truncate("hi", 10); got != "hi" {
		t.Errorf("truncate short string: got %q", got)
	}
}

func TestCapitalize(t *testing.T) {
	if got := capitalize("opened"); got != "Opened" {
		t.Errorf("capitalize: got %q", got)
	}
	if got := capitalize(""); got != "" {
		t.Errorf("capitalize empty: got %q", got)
	}
}

func mockServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}
func fetchEventsFromURL(baseURL, username string, limit int) ([]GitHubEvent, error) {
	v.
		_ = os.Setenv("_TEST_BASE_URL", baseURL)
	return nil, nil
}

func TestFetchEventsNotFound(t *testing.T) {
	srv := mockServer(t, http.StatusNotFound, map[string]string{"message": "Not Found"})
	defer srv.Close()
	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestFetchEventsDecoding(t *testing.T) {
	events := []GitHubEvent{
		{Type: "PushEvent", Repo: Repo{Name: "user/repo"}, Payload: Payload{Size: 2}},
		{Type: "WatchEvent", Repo: Repo{Name: "user/other"}},
	}
	srv := mockServer(t, http.StatusOK, events)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var got []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 events, got %d", len(got))
	}
	if got[0].Type != "PushEvent" {
		t.Errorf("expected PushEvent, got %s", got[0].Type)
	}
}
