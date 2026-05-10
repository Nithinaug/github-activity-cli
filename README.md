# GitHub Activity CLI

[![CI](https://github.com/Nithinaug/github-activity-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Nithinaug/github-activity-cli/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A fast, colorful CLI tool written in **Go** that displays a developer's recent GitHub activity right in your terminal — pushes, PRs, issues, stars, forks, and more.

```
Fetching activity for torvalds...

  ↑  Pushed 4 commits to torvalds/linux
  ↑  Pushed 1 commit to torvalds/linux
  ★  Starred git/git
  ⇄  Opened pull request #1234 in torvalds/linux  Add new scheduler option
  ⑂  Forked golang/go → torvalds/go

── Activity Summary ──────────────────
  Pushes                         12
  Stars Given                     3
  Pull Requests                   1
  Forks                           1
──────────────────────────────────────
```

---

## Features

- **15+ event types** — pushes, PRs, issues, comments, stars, forks, releases, and more
- **Colored output** — green for pushes, yellow for issues, cyan for creates, no external dependencies
- **`--filter` flag** — focus on just push, issue, or PR events
- **`--limit` flag** — control how many events to show
- **Authenticated requests** via `GITHUB_TOKEN` (5 000 req/hr vs 60 unauthenticated)
- **Informative errors** — clear messages for 404s and rate limits
- **Unit tested** with `net/http/httptest` mocks
- **Cross-platform** — Linux, macOS, Windows binaries built by CI

---

## Installation

### Install with Go

```bash
go install github.com/Nithinaug/github-activity-cli@latest
```

### Download a binary

Pre-built binaries are available on the [Releases](https://github.com/Nithinaug/github-activity-cli/releases) page.

### Build from source

```bash
git clone https://github.com/Nithinaug/github-activity-cli.git
cd github-activity-cli
go build -o github-activity .
```

---

## Usage

```
github-activity <username> [flags]

Flags:
  -filter string     Show only events whose type contains this string (e.g. push, issue, pr)
  -limit int         Maximum number of events to display, 0 = all (default 30)
  -no-summary        Skip the aggregated summary table
  -version           Print version and exit

Environment:
  GITHUB_TOKEN       Personal access token for authenticated requests (5000 req/hr)
```

### Examples

```bash
# Basic usage
github-activity torvalds

# Show only push events
github-activity torvalds --filter push

# Show last 10 events, no summary
github-activity torvalds --limit 10 --no-summary

# Authenticated (higher rate limit)
export GITHUB_TOKEN=ghp_yourtoken
github-activity torvalds
```

---

## Tech Stack

| | |
|---|---|
| Language | Go 1.22 |
| API | GitHub REST API v2022-11-28 |
| Testing | `testing` + `net/http/httptest` |
| CI | GitHub Actions |
| Dependencies | **Zero** — stdlib only |

---

## Project Structure

```
.
├── main.go          # CLI entry point, flag parsing
├── github.go        # GitHub API client (FetchEvents)
├── summary.go       # Event formatting and colored output
├── types.go         # Data structures for all GitHub event types
├── activity_test.go # Unit tests
└── .github/
    └── workflows/
        └── ci.yml   # Build, test, cross-compile on every push
```

---

## License

MIT