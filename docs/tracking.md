# Orbit Tracking Command Reference

Push workflow events and token usage data to Draxarp Tracking for analytics and compliance reporting.

Requires a profile with a `draxarp` service configured.

---

## event push

Push a workflow event to Draxarp Tracking.

```
orbit tracking event push [flags] -p myprofile
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--type` | string | Event type (required). E.g., `skill.completed`, `jira.transition`, `agent.invoked`. |
| `--ticket` | string | Jira ticket ID (e.g., `FOUN-42`). |
| `--project-key` | string | Jira project key (e.g., `FOUN`). |
| `--skill` | string | Skill name (e.g., `execute`, `commit`). |
| `--agent` | string | Agent name (e.g., `developer`, `architect`). |
| `--phase-from` | string | PIV phase transitioning from. |
| `--phase-to` | string | PIV phase transitioning to. |
| `--status-from` | string | Jira status transitioning from. |
| `--status-to` | string | Jira status transitioning to. |
| `--session` | string | Claude Code session ID. |
| `--context-reset` | int | Context reset number. |
| `--meta` | string | JSON metadata (e.g., `'{"branch":"feature/x"}'`). |

**Examples:**

```bash
# Skill completed
orbit -p paybook tracking event push --type "skill.completed" --ticket FOUN-42 --skill execute

# Jira transition
orbit -p paybook tracking event push --type "jira.transition" --ticket FOUN-42 --status-to "In Code Review"

# Agent invoked
orbit -p paybook tracking event push --type "agent.invoked" --ticket FOUN-42 --agent developer
```

---

## tkm push

Push unpushed token usage events from the local TKM database to Draxarp Tracking. Uses incremental sync — only pushes events not yet sent.

```
orbit tracking tkm push [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--batch-size` | int | 100 | Max events to push per run. |
| `--dry-run` | bool | false | Show events without pushing. |

**Examples:**

```bash
# Push unpushed events
orbit -p paybook tracking tkm push

# Push with larger batch
orbit -p paybook tracking tkm push --batch-size 500

# Preview without pushing
orbit -p paybook tracking tkm push --dry-run
```

---

## Invocation Source & Tool Version Detection

Workflow events are automatically tagged with their invocation source and tool versions:

| Context | `metadata.invoked_by` | `metadata.entrypoint` |
|---------|----------------------|-----------------------|
| Claude Code CLI session | `claude-code` | `cli` |
| Claude Code IDE extension | `claude-code` | (varies) |
| User running orbit manually | `manual` | — |

### Tool Versions (transparent)

Every event also includes `metadata.tool_versions` with:

| Key | Example | Source |
|-----|---------|--------|
| `orbit` | `0.46.0` | `orbit version` |
| `claude_code` | `1.0.45` | `claude --version` |
| `paybook_workflow` | `1.2.0` | Plugin list |

These are captured transparently — no user input needed. Draxarp uses them to track tool adoption rates and version distribution across the team.

## How It Works

1. `orbit tkm sync` ingests Claude Code session data into `~/.config/orbit/tkm.db`
2. `orbit tracking tkm push` reads unpushed events from the local DB and sends them to Draxarp
3. A `push_state` table tracks the last pushed event ID per target — subsequent pushes only send new events
4. Events are sent in chunks (default 100) to the Draxarp `/api/v1/tracking/ingest/token-usage` endpoint
5. Draxarp deduplicates by `event_hash` — re-pushing is safe
