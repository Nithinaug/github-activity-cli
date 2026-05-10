package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "2.0.0"

func main() {
	limit := flag.Int("limit", 30, "Maximum number of events to display (0 = all)")
	filter := flag.String("filter", "", "Show only events whose type contains this string (e.g. push, issue, pr)")
	noSummary := flag.Bool("no-summary", false, "Skip the aggregated summary table")
	showVer := flag.Bool("version", false, "Print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: github-activity <username> [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment:\n")
		fmt.Fprintf(os.Stderr, "  GITHUB_TOKEN  Personal access token for authenticated requests (5000 req/hr)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  github-activity torvalds\n")
		fmt.Fprintf(os.Stderr, "  github-activity torvalds --filter push --limit 10\n")
		fmt.Fprintf(os.Stderr, "  github-activity torvalds --no-summary\n")
	}

	flag.Parse()

	if *showVer {
		fmt.Println("github-activity version", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: username is required")
		flag.Usage()
		os.Exit(1)
	}

	username := flag.Arg(0)

	fmt.Printf("Fetching activity for %s...\n\n", username)

	events, err := FetchEvents(username, *limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No public activity found.")
		os.Exit(0)
	}

	PrintEvents(events, *filter)

	if !*noSummary {
		PrintSummary(events)
	}
}
