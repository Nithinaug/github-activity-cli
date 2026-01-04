package main

import "fmt"

func printSummary(events []GitHubEvent) {
	counts := make(map[string]int)
    seen := make(map[string]bool)
for _, event := range events {
    counts[event.Type]++
    switch event.Type {
    case "PushEvent":
        if !seen[event.Repo.Name] {
            fmt.Println("Pushed to", event.Repo.Name)
            seen[event.Repo.Name] = true
        }
    case "CreateEvent":
        fmt.Println("Created", event.Repo.Name)
    case "WatchEvent":
        fmt.Println("Starred", event.Repo.Name)
    case "PullRequestEvent":
        fmt.Println("Pull request activity in", event.Repo.Name)
    }
}
fmt.Println("--------------------")
fmt.Println("Summary")
fmt.Println("--------------------")
for eventType, count := range counts {
    fmt.Printf("%-20s %d\n", eventType, count)
}
}