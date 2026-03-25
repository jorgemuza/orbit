---
name: tkm
description: "Track AI token consumption, costs, and usage trends using the orbit CLI. Use this skill whenever the user asks about token usage, AI costs, Claude Code spending, how many tokens were used, cost breakdown by model, session history, or token analytics. Trigger on phrases like 'how much have I spent', 'token usage', 'show me costs', 'what's my AI spending', 'how many tokens today', 'cost per model', 'list sessions', 'track usage', 'token report', 'weekly usage', 'monthly costs', or any token/cost tracking task — even casual references like 'am I spending too much on Claude', 'what did that session cost', 'show me the dashboard', or 'how much is opus costing us'."
---

# Token Manager (tkm) with orbit CLI

Track AI token consumption, costs, and usage trends across Claude Code sessions. Data is ingested automatically from Claude Code session JSONL files and stored in a local SQLite database.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. Claude Code installed and used (session files in `~/.claude/projects/`)

No profile or service configuration needed — tkm works entirely with local data.

## Quick Reference

```bash
# Today's token usage and cost
orbit tkm status

# This week's summary
orbit tkm status --period week

# Daily usage breakdown
orbit tkm usage

# Weekly usage
orbit tkm usage --period weekly

# Monthly usage
orbit tkm usage --period monthly

# Cost breakdown by model
orbit tkm cost

# Monthly cost trends
orbit tkm cost --period month

# Recent sessions with token counts
orbit tkm sessions

# More sessions
orbit tkm sessions --limit 50

# Filter by project
orbit tkm status --project /path/to/project
orbit tkm usage --project /path/to/project

# Filter by model
orbit tkm usage --model opus

# Force sync session data
orbit tkm sync

# Full re-sync (re-process all files)
orbit tkm sync --full

# JSON output for scripting
orbit tkm status -o json
orbit tkm cost -o json

# Manually record a usage event
orbit tkm track --model claude-opus-4-6 --input 1000 --output 500

# Clear all tracking data
orbit tkm clear --confirm
```

## Core Workflows

### Check daily spending

```bash
orbit tkm status
orbit tkm cost --period today
```

### Weekly cost review

```bash
orbit tkm usage --period weekly --limit 8
orbit tkm cost --period week
```

### Per-project analysis

```bash
# Check spending on a specific project
orbit tkm status --period month --project /Users/me/Projects/my-app
orbit tkm sessions --project /Users/me/Projects/my-app
```

### Model cost comparison

```bash
orbit tkm cost --period month
```

### Export data for reporting

```bash
orbit tkm usage --period monthly -o json > monthly-usage.json
orbit tkm cost -o json > cost-breakdown.json
orbit tkm sessions -o json > sessions.json
```

## Data Storage

- **Database:** `~/.config/orbit/tkm.db` (SQLite)
- **Source:** `~/.claude/projects/<encoded-path>/<session-id>.jsonl`
- **Auto-sync:** Runs automatically on every read command (status, usage, cost, sessions)
- **Incremental:** Only parses new data since last sync (byte-offset tracking)
- **Retention:** Skips files older than 90 days

## Model Pricing

Default pricing (per million tokens):

| Model | Input | Output | Cache Read | Cache Write |
|-------|-------|--------|------------|-------------|
| claude-opus-4-6 | $15.00 | $75.00 | $1.50 | $18.75 |
| claude-sonnet-4-6 | $3.00 | $15.00 | $0.30 | $3.75 |
| claude-haiku-4-5 | $0.80 | $4.00 | $0.08 | $1.00 |

For full command details and flags, see `references/commands.md`.
