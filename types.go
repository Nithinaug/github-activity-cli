package main

type GitHubEvent struct {
	Type string `json:"type"`

	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`

	CreatedAt string `json:"created_at"`
}
