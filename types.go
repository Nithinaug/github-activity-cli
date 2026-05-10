package main

type GitHubEvent struct {
	Type      string  `json:"type"`
	Repo      Repo    `json:"repo"`
	Payload   Payload `json:"payload"`
	CreatedAt string  `json:"created_at"`
}

type Repo struct {
	Name string `json:"name"`
}
type Payload struct {
	// PushEvent
	Commits []Commit `json:"commits"`
	Size    int      `json:"size"`
	Action  string   `json:"action"`

	Issue *Issue `json:"issue"`

	PullRequest *PullRequest `json:"pull_request"`
	RefType     string       `json:"ref_type"`
	Ref         string       `json:"ref"`

	Forkee  *Forkee  `json:"forkee"`
	Comment *Comment `json:"comment"`
}

type Commit struct {
	Message string `json:"message"`
}

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

type Forkee struct {
	FullName string `json:"full_name"`
}

type Comment struct {
	Body string `json:"body"`
}
