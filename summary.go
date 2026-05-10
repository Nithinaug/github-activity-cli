package main

import (
	"fmt"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBlue   = "\033[34m"
	colorRed    = "\033[31m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"
	bold        = "\033[1m"
)

func color(c, s string) string {
	return c + s + colorReset
}

func PrintEvents(events []GitHubEvent, filter string) {
	filter = strings.ToLower(filter)
	printed := 0

	for _, e := range events {
		if filter != "" && !strings.Contains(strings.ToLower(e.Type), filter) {
			continue
		}
		line := formatEvent(e)
		if line != "" {
			fmt.Println(line)
			printed++
		}
	}

	if printed == 0 {
		fmt.Println(color(colorGray, "  No matching events found."))
	}
}

func PrintSummary(events []GitHubEvent) {
	counts := make(map[string]int)
	for _, e := range events {
		counts[e.Type]++
	}

	fmt.Println()
	fmt.Println(color(bold+colorWhite, "── Activity Summary ──────────────────"))
	for _, et := range eventOrder {
		if n, ok := counts[et]; ok {
			fmt.Printf("  %-30s %s\n", color(colorGray, friendlyName(et)), color(colorCyan, fmt.Sprintf("%d", n)))
			delete(counts, et)
		}
	}
	for et, n := range counts {
		fmt.Printf("  %-30s %s\n", color(colorGray, friendlyName(et)), color(colorCyan, fmt.Sprintf("%d", n)))
	}
	fmt.Println(color(colorGray, "──────────────────────────────────────"))
}

func formatEvent(e GitHubEvent) string {
	repo := color(colorBlue, e.Repo.Name)

	switch e.Type {
	case "PushEvent":
		n := e.Payload.Size
		if n == 0 {
			n = len(e.Payload.Commits)
		}
		commits := "commit"
		if n != 1 {
			commits = "commits"
		}
		return fmt.Sprintf("  %s  Pushed %s %s to %s",
			color(colorGreen, "↑"),
			color(colorGreen, fmt.Sprintf("%d", n)),
			commits, repo)

	case "CreateEvent":
		ref := e.Payload.Ref
		if ref == "" {
			ref = e.Payload.RefType
		}
		return fmt.Sprintf("  %s  Created %s %s in %s",
			color(colorCyan, "✦"),
			color(colorGray, e.Payload.RefType),
			color(colorCyan, ref), repo)

	case "DeleteEvent":
		return fmt.Sprintf("  %s  Deleted %s %s in %s",
			color(colorRed, "✗"),
			color(colorGray, e.Payload.RefType),
			color(colorRed, e.Payload.Ref), repo)

	case "WatchEvent":
		return fmt.Sprintf("  %s  Starred %s", color(colorYellow, "★"), repo)

	case "ForkEvent":
		fork := repo
		if e.Payload.Forkee != nil {
			fork = color(colorBlue, e.Payload.Forkee.FullName)
		}
		return fmt.Sprintf("  %s  Forked %s → %s", color(colorCyan, "⑂"), repo, fork)

	case "IssuesEvent":
		action := e.Payload.Action
		num, title := "", ""
		if e.Payload.Issue != nil {
			num = fmt.Sprintf("#%d", e.Payload.Issue.Number)
			title = color(colorGray, truncate(e.Payload.Issue.Title, 50))
		}
		return fmt.Sprintf("  %s  %s issue %s in %s %s",
			color(colorYellow, "⚑"),
			capitalize(action),
			color(colorYellow, num), repo, title)

	case "PullRequestEvent":
		action := e.Payload.Action
		num, title := "", ""
		if e.Payload.PullRequest != nil {
			num = fmt.Sprintf("#%d", e.Payload.PullRequest.Number)
			title = color(colorGray, truncate(e.Payload.PullRequest.Title, 50))
		}
		return fmt.Sprintf("  %s  %s pull request %s in %s %s",
			color(colorGreen, "⇄"),
			capitalize(action),
			color(colorGreen, num), repo, title)

	case "IssueCommentEvent":
		num := ""
		if e.Payload.Issue != nil {
			num = fmt.Sprintf("#%d", e.Payload.Issue.Number)
		}
		return fmt.Sprintf("  %s  Commented on issue %s in %s",
			color(colorGray, "💬"),
			color(colorYellow, num), repo)

	case "PullRequestReviewEvent":
		return fmt.Sprintf("  %s  Reviewed a pull request in %s",
			color(colorGreen, "✔"), repo)

	case "PullRequestReviewCommentEvent":
		return fmt.Sprintf("  %s  Commented on a PR review in %s",
			color(colorGray, "💬"), repo)

	case "ReleaseEvent":
		return fmt.Sprintf("  %s  Published a release in %s",
			color(colorCyan, "🚀"), repo)

	case "PublicEvent":
		return fmt.Sprintf("  %s  Made %s public",
			color(colorGreen, "🌐"), repo)

	case "MemberEvent":
		return fmt.Sprintf("  %s  %s a collaborator in %s",
			color(colorBlue, "👤"),
			capitalize(e.Payload.Action), repo)

	case "GollumEvent":
		return fmt.Sprintf("  %s  Updated wiki in %s",
			color(colorGray, "📄"), repo)

	case "CommitCommentEvent":
		return fmt.Sprintf("  %s  Commented on a commit in %s",
			color(colorGray, "💬"), repo)

	default:
		return fmt.Sprintf("  %s  %s in %s",
			color(colorGray, "·"),
			color(colorGray, e.Type), repo)
	}
}

var eventOrder = []string{
	"PushEvent", "PullRequestEvent", "IssuesEvent", "IssueCommentEvent",
	"PullRequestReviewEvent", "PullRequestReviewCommentEvent",
	"CreateEvent", "DeleteEvent", "WatchEvent", "ForkEvent",
	"ReleaseEvent", "PublicEvent", "MemberEvent", "GollumEvent", "CommitCommentEvent",
}

func friendlyName(t string) string {
	names := map[string]string{
		"PushEvent":                     "Pushes",
		"PullRequestEvent":              "Pull Requests",
		"IssuesEvent":                   "Issues",
		"IssueCommentEvent":             "Issue Comments",
		"PullRequestReviewEvent":        "PR Reviews",
		"PullRequestReviewCommentEvent": "PR Review Comments",
		"CreateEvent":                   "Branches/Tags Created",
		"DeleteEvent":                   "Branches/Tags Deleted",
		"WatchEvent":                    "Stars Given",
		"ForkEvent":                     "Forks",
		"ReleaseEvent":                  "Releases",
		"PublicEvent":                   "Repos Made Public",
		"MemberEvent":                   "Collaborator Events",
		"GollumEvent":                   "Wiki Updates",
		"CommitCommentEvent":            "Commit Comments",
	}
	if n, ok := names[t]; ok {
		return n
	}
	return t
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
