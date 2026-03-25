# tkm Command Reference

Complete reference for all `orbit tkm` commands and flags.

## status

Quick summary of token usage and cost for a time period.

| Flag | Default | Description |
|------|---------|-------------|
| `--period` | `today` | Time period: `today`, `week`, `month`, `all` |
| `--project` | | Filter by project path |

```
orbit tkm status
orbit tkm status --period week
orbit tkm status --project /path/to/project
```

---

## usage

Detailed token usage report with time aggregation.

| Flag | Default | Description |
|------|---------|-------------|
| `--period` | `daily` | Aggregation: `daily`, `weekly`, `monthly` |
| `--project` | | Filter by project path |
| `--model` | | Filter by model name (partial match) |
| `--limit` | 30 | Max rows |

```
orbit tkm usage
orbit tkm usage --period weekly
orbit tkm usage --period monthly --model opus
```

---

## cost

Cost breakdown by model.

| Flag | Default | Description |
|------|---------|-------------|
| `--period` | `all` | Time period: `today`, `week`, `month`, `all` |
| `--project` | | Filter by project path |

```
orbit tkm cost
orbit tkm cost --period month
```

---

## sessions

List recent sessions with aggregated token counts and costs.

| Flag | Default | Description |
|------|---------|-------------|
| `--project` | | Filter by project path |
| `--limit` | 20 | Max sessions |

```
orbit tkm sessions
orbit tkm sessions --limit 50
```

---

## sync

Force ingestion of Claude Code session data.

| Flag | Default | Description |
|------|---------|-------------|
| `--full` | false | Re-process all files from scratch |

```
orbit tkm sync
orbit tkm sync --full
```

---

## track

Manually record a token usage event.

| Flag | Default | Description |
|------|---------|-------------|
| `--model` | | Model name (required) |
| `--input` | 0 | Input tokens |
| `--output` | 0 | Output tokens |
| `--cache-read` | 0 | Cache read tokens |
| `--cache-write` | 0 | Cache creation tokens |
| `--session-id` | `manual` | Session ID |
| `--project` | cwd | Project path |

```
orbit tkm track --model claude-opus-4-6 --input 1000 --output 500
```

---

## clear

Clear tracking data.

| Flag | Default | Description |
|------|---------|-------------|
| `--confirm` | false | Required to confirm deletion |
| `--project` | | Clear data for specific project only |

```
orbit tkm clear --confirm
orbit tkm clear --confirm --project /path/to/project
```
