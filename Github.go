package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func FetchEvents(username string, limit int) ([]GitHubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public?per_page=100", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gh-activity-cli")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("user '%s' not found", username)
	case http.StatusForbidden, http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded — set GITHUB_TOKEN for higher limits")
	default:
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var events []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	if limit > 0 && limit < len(events) {
		events = events[:limit]
	}

	return events, nil
}
