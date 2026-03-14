# Orbit Git Command Reference

Local git operations from the orbit CLI.

**Top-level command:** `orbit git`

**Notes:**

- These commands wrap local `git` operations. No service connection or profile is required.
- Unlike `orbit github` / `orbit gitlab`, these operate on the local repository.

---

## Table of Contents

- [push — Push commits](#push)
- [branch — Branch commands](#branch)
  - [branch list](#branch-list)
  - [branch create](#branch-create)
  - [branch delete](#branch-delete)
  - [branch current](#branch-current)
  - [branch switch](#branch-switch)

---

## push

Push commits to a remote repository.

```
orbit git push [remote] [branch] [flags]
```

| Argument | Default | Description |
|----------|---------|-------------|
| `remote` | `origin` | Remote name |
| `branch` | current branch | Branch to push |

| Flag | Short | Description |
|------|-------|-------------|
| `--set-upstream` | `-u` | Set upstream tracking reference |
| `--force` | `-f` | Force push |

**Examples:**

```bash
# Push current branch to origin
orbit git push

# Push with upstream tracking
orbit git push -u origin

# Push a specific branch
orbit git push origin main

# Push and set upstream for a feature branch
orbit git push -u origin feature/my-branch
```

---

## branch

Manage local git branches.

### branch list

List branches.

```
orbit git branch list [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | `-a` | List local and remote branches |
| `--remote` | `-r` | List remote branches only |

**Examples:**

```bash
# List local branches
orbit git branch list

# List all branches (local + remote)
orbit git branch list --all

# List remote branches only
orbit git branch list --remote
```

### branch create

Create a new branch.

```
orbit git branch create [name] [start-point]
```

| Argument | Required | Description |
|----------|----------|-------------|
| `name` | Yes | Branch name |
| `start-point` | No | Starting commit or branch (defaults to HEAD) |

**Examples:**

```bash
# Create a branch from HEAD
orbit git branch create feature/my-feature

# Create a branch from main
orbit git branch create bugfix/fix-123 main
```

### branch delete

Delete a branch.

```
orbit git branch delete [name] [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Branch name to delete |

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Force delete (even if not merged) |

**Examples:**

```bash
# Delete a merged branch
orbit git branch delete feature/old-feature

# Force delete an unmerged branch
orbit git branch delete -f feature/unmerged
```

### branch current

Show the current branch name.

```
orbit git branch current
```

**Example:**

```bash
orbit git branch current
```

### branch switch

Switch to a branch.

```
orbit git branch switch [name] [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Branch name to switch to |

| Flag | Short | Description |
|------|-------|-------------|
| `--create` | `-c` | Create the branch if it doesn't exist |

**Examples:**

```bash
# Switch to an existing branch
orbit git branch switch main

# Create and switch to a new branch
orbit git branch switch -c feature/new-feature
```
