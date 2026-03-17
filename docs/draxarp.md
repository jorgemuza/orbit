# orbit draxarp

Manage Draxarp Development Intelligence — projects, tasks, specs, docs, and memories.

**Alias:** `dx`

## Global Flags

| Flag | Description |
|------|-------------|
| `--service <name>` | Select draxarp service (if profile has multiple) |
| `-o json` | JSON output |
| `-p <profile>` | Use specific profile |

## Configuration

```yaml
# ~/.config/orbit/config.yaml
profiles:
  - name: my-project
    services:
      - name: draxarp
        type: draxarp
        base_url: https://your-instance.com
        auth:
          method: token
          token: "your-api-token"
```

---

## Projects

Intelligence projects group tasks, specs, memories, and docs.

```bash
# List projects
orbit draxarp project list
orbit dx proj ls

# View a project
orbit draxarp project view <project-id>

# Create a project
orbit draxarp project create --name "My Project" --description "..." --repo https://github.com/...

# Delete a project
orbit draxarp project delete <project-id>
```

**Flags for `project list`:**
| Flag | Description |
|------|-------------|
| `--limit <n>` | Max results (default: 20) |

**Flags for `project create`:**
| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Project name |
| `--description` | No | Description |
| `--repo` | No | Repository URL |

---

## Tasks

Development tasks with status tracking and priorities.

```bash
# List tasks
orbit draxarp task list
orbit dx task ls --project <id> --status in_progress

# View a task
orbit draxarp task view <task-id>

# Create a task
orbit draxarp task create --title "Fix auth bug" --project <id> --priority high

# Complete a task
orbit draxarp task complete <task-id>

# Delete a task
orbit draxarp task delete <task-id>
```

**Flags for `task list`:**
| Flag | Description |
|------|-------------|
| `--project <id>` | Filter by project |
| `--status <s>` | Filter: `pending`, `in_progress`, `completed`, `blocked`, `cancelled` |
| `--limit <n>` | Max results (default: 20) |

**Flags for `task create`:**
| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Task title |
| `--project` | Yes | Project ID |
| `--description` | No | Description |
| `--priority` | No | `low`, `medium`, `high`, `critical` |
| `--status` | No | Initial status |

---

## Memories

Persistent knowledge store — patterns, bugs, decisions, domain context.

```bash
# List memories
orbit draxarp memory list
orbit dx mem ls --project <id> --type architecture

# View a memory
orbit draxarp memory view <memory-id>

# Create a memory
orbit draxarp memory create --title "Pattern: Repository" --content "..." --project <id> --type patterns

# Archive a memory
orbit draxarp memory archive <memory-id>
```

**Memory types:** `architecture`, `patterns`, `bugs`, `domain`, `decisions`, `preferences`, `context`

**Flags for `memory list`:**
| Flag | Description |
|------|-------------|
| `--project <id>` | Filter by project |
| `--type <t>` | Filter by memory type |
| `--status <s>` | Filter by status |
| `--limit <n>` | Max results (default: 20) |

**Flags for `memory create`:**
| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Memory title |
| `--content` | Yes | Memory content |
| `--project` | Yes | Project ID |
| `--type` | Yes | Memory type (see above) |

---

## Specs

Specifications with review workflow: draft → in_review → approved/rejected.

```bash
# List specs
orbit draxarp spec list
orbit dx spec ls --project <id> --status draft

# View a spec
orbit draxarp spec view <spec-id>

# Create a spec
orbit draxarp spec create --title "Auth Spec" --content "..." --project <id>

# Submit for review
orbit draxarp spec submit <spec-id>

# Approve / Reject
orbit draxarp spec approve <spec-id>
orbit draxarp spec reject <spec-id>

# Delete a spec
orbit draxarp spec delete <spec-id>
```

**Flags for `spec list`:**
| Flag | Description |
|------|-------------|
| `--project <id>` | Filter by project |
| `--status <s>` | Filter: `draft`, `in_review`, `approved`, `rejected` |
| `--limit <n>` | Max results (default: 20) |

**Flags for `spec create`:**
| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Spec title |
| `--content` | Yes | Spec content |
| `--project` | Yes | Project ID |

---

## Docs

Technical documentation with slug-based addressing and wikilink support.

```bash
# List docs
orbit draxarp doc list
orbit dx doc ls --category architecture --status published

# Table of contents (lightweight — no content)
orbit draxarp doc toc

# View a doc by ID
orbit draxarp doc view <doc-id>

# View a doc by slug
orbit draxarp doc view architecture/overview --slug

# Create a doc
orbit draxarp doc create --title "API Guide" --category api --slug api/guide

# Publish / Archive
orbit draxarp doc publish <doc-id>
orbit draxarp doc archive <doc-id>

# Delete
orbit draxarp doc delete <doc-id>
```

**Flags for `doc list`:**
| Flag | Description |
|------|-------------|
| `--project <id>` | Filter by project |
| `--status <s>` | Filter: `draft`, `published`, `archived` |
| `--category <c>` | Filter by category (see below) |
| `--limit <n>` | Max results (default: 20) |

**Doc categories:** `architecture`, `components`, `flows`, `security`, `api`, `infrastructure`, `general`

**Flags for `doc create`:**
| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Doc title |
| `--content` | No | Doc content |
| `--project` | No | Project ID |
| `--category` | No | Category |
| `--slug` | No | Custom slug (auto-generated if empty) |

**Flags for `doc view`:**
| Flag | Description |
|------|-------------|
| `--slug` | Look up by slug instead of ID |
