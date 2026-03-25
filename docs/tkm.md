# Orbit Token Manager (tkm) Command Reference

Track AI token consumption, costs, and usage trends across Claude Code sessions.

Data is ingested automatically from Claude Code session files (`~/.claude/projects/`). All commands support `-o json` and `-o yaml` output formats.

---

## status

Quick summary of token usage and cost.

```
orbit tkm status [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--period` | string | `today` | Time period: `today`, `week`, `month`, `all`. |
| `--project` | string | | Filter by project path. |

**Examples:**

```bash
orbit tkm status
orbit tkm status --period week
orbit tkm status --project /path/to/project
orbit tkm status -o json
```

---

## usage

Detailed token usage report with time aggregation.

```
orbit tkm usage [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--period` | string | `daily` | Aggregation: `daily`, `weekly`, `monthly`. |
| `--project` | string | | Filter by project path. |
| `--model` | string | | Filter by model name (partial match). |
| `--limit` | int | 30 | Max rows to show. |

**Examples:**

```bash
orbit tkm usage
orbit tkm usage --period weekly
orbit tkm usage --period monthly --model opus
orbit tkm usage --project /path/to/project --limit 10
```

---

## cost

Cost breakdown by model.

```
orbit tkm cost [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--period` | string | `all` | Time period: `today`, `week`, `month`, `all`. |
| `--project` | string | | Filter by project path. |

**Examples:**

```bash
orbit tkm cost
orbit tkm cost --period month
orbit tkm cost --project /path/to/project
```

---

## sessions

List recent Claude Code sessions with token counts and costs.

```
orbit tkm sessions [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | | Filter by project path. |
| `--limit` | int | 20 | Max sessions to show. |

**Examples:**

```bash
orbit tkm sessions
orbit tkm sessions --limit 50
orbit tkm sessions --project /path/to/project
```

---

## sync

Force ingestion of Claude Code session data.

```
orbit tkm sync [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--full` | bool | false | Re-process all files from scratch (ignores saved offsets). |

**Examples:**

```bash
orbit tkm sync
orbit tkm sync --full
```

---

## track

Manually record a token usage event. Useful for hook integration or tracking non-Claude tools.

```
orbit tkm track [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--model` | string | | Model name (required). |
| `--input` | int | 0 | Input tokens. |
| `--output` | int | 0 | Output tokens. |
| `--cache-read` | int | 0 | Cache read tokens. |
| `--cache-write` | int | 0 | Cache creation tokens. |
| `--session-id` | string | `manual` | Session ID. |
| `--project` | string | cwd | Project path. |

**Examples:**

```bash
orbit tkm track --model claude-opus-4-6 --input 1000 --output 500
orbit tkm track --model claude-sonnet-4-6 --input 5000 --output 2000 --cache-read 10000
```

---

## clear

Clear tracking data.

```
orbit tkm clear [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--confirm` | bool | false | Required to confirm deletion. |
| `--project` | string | | Clear data for specific project only. |

**Examples:**

```bash
orbit tkm clear --confirm
orbit tkm clear --confirm --project /path/to/project
```

---

## How It Works

1. Claude Code stores session data as JSONL files in `~/.claude/projects/<encoded-path>/<session-id>.jsonl`
2. On any read command (`status`, `usage`, `cost`, `sessions`), tkm auto-syncs new events
3. Each assistant response with token usage is extracted, cost is calculated using per-model pricing, and stored in a local SQLite database at `~/.config/orbit/tkm.db`
4. Incremental parsing uses byte-offset tracking to avoid re-reading entire files

## Default Model Pricing (per million tokens)

| Model | Input | Output | Cache Read | Cache Write |
|-------|-------|--------|------------|-------------|
| claude-opus-4-6 | $15.00 | $75.00 | $1.50 | $18.75 |
| claude-sonnet-4-6 | $3.00 | $15.00 | $0.30 | $3.75 |
| claude-haiku-4-5 | $0.80 | $4.00 | $0.08 | $1.00 |
