# GitLab CLI Command Reference

Complete reference for all `orbit gitlab` (alias: `gl`) commands and flags.

## Table of Contents

- [project / projects](#project--projects)
- [group](#group)
- [branch](#branch)
- [tag](#tag)
- [commit](#commit)
- [mr (merge request)](#mr-merge-request)
- [pipeline (aliases: pipe, ci)](#pipeline)
- [schedule (alias: sched)](#schedule)
- [issue](#issue)
- [member](#member)
- [user](#user)

## Global Flags

These flags apply to all gitlab commands:

| Flag | Description |
|------|-------------|
| `-p, --profile <name>` | Profile to use (required) |
| `-o, --output <format>` | Output format: `table` (default), `json`, `yaml` |
| `--service <name>` | GitLab service name (if profile has multiple) |

---

## project / projects

### `gitlab project <id-or-path>`

View a single project.

```
orbit -p myprofile gl project 595
orbit -p myprofile gl project schools/frontend/my-app
```

**Output fields:** ID, Name, Path, Description, Default Branch, Visibility, URL, Last Activity

### `gitlab project list`

List projects.

| Flag | Default | Description |
|------|---------|-------------|
| `--search <text>` | | Search by name |
| `--group <id-or-path>` | | List projects in a group (includes subgroups) |
| `--limit <n>` | 50 | Max results |

```
orbit -p myprofile gl project list --search frontend
orbit -p myprofile gl project list --group schools/frontend
```

### `gitlab project activity`

List projects with recent activity, showing which branches received pushes.

| Flag | Default | Description |
|------|---------|-------------|
| `--since <YYYY-MM-DD>` | | Show activity since date |
| `--days <n>` | | Show activity from the last N days |
| `--group <id-or-path>` | | Filter to projects in a group |
| `--branch <names>` | | Filter by branch names (comma-separated) |
| `--limit <n>` | 100 | Max projects to scan |

Defaults to today when neither `--since` nor `--days` is provided.

```
orbit -p myprofile gl project activity
orbit -p myprofile gl project activity --days 7
orbit -p myprofile gl project activity --group schools --branch development,master
```

---

## group

### `gitlab group view <id-or-path>`

View group details.

```
orbit -p myprofile gl group view schools/frontend
```

### `gitlab group list`

List groups.

| Flag | Default | Description |
|------|---------|-------------|
| `--search <text>` | | Search by name |
| `--limit <n>` | 50 | Max results |

### `gitlab group subgroups <id-or-path>`

List subgroups of a group.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

---

## branch

### `gitlab branch list <project>`

List branches.

| Flag | Default | Description |
|------|---------|-------------|
| `--search <text>` | | Search by name |
| `--limit <n>` | 50 | Max results |

### `gitlab branch view <project> <branch>`

View branch details including latest commit.

### `gitlab branch create <project> <name> <ref>`

Create a branch from a ref (branch name, tag, or SHA).

```
orbit -p myprofile gl branch create 595 feature/new-thing main
```

### `gitlab branch delete <project> <branch>`

Delete a branch.

```
orbit -p myprofile gl branch delete 595 feature/old-thing
```

### `gitlab branch protect <project> <branch>`

Protect a branch with push, merge, and unprotect access controls.

| Flag | Default | Description |
|------|---------|-------------|
| `--push` | `maintainer` | Push access: `no-access`, `developer`, `maintainer`, `admin` |
| `--merge` | `maintainer` | Merge access: `no-access`, `developer`, `maintainer`, `admin` |
| `--unprotect` | | Unprotect access level |
| `--allow-force-push` | `false` | Allow force push |
| `--allowed-to-push` | | User IDs allowed to push (repeatable) |
| `--allowed-to-merge` | | User IDs allowed to merge (repeatable) |

```
orbit -p myprofile gl branch protect 650 main --push no-access --merge maintainer
orbit -p myprofile gl branch protect 650 main --push no-access --merge maintainer --allowed-to-push 12
```

### `gitlab branch unprotect <project> <branch>`

Remove branch protection.

```
orbit -p myprofile gl branch unprotect 650 main
```

### `gitlab branch protections <project>`

List all protected branches with push/merge access levels.

```
orbit -p myprofile gl branch protections 650
```

---

## tag

### `gitlab tag list <project>`

List tags.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `gitlab tag create <project> <tag-name> <ref>`

Create a tag from a ref.

| Flag | Description |
|------|-------------|
| `-m, --message <text>` | Tag message (annotated tag) |

```
orbit -p myprofile gl tag create 595 v1.0.0 main -m "Release v1.0.0"
```

---

## commit

### `gitlab commit list <project>`

List recent commits.

| Flag | Default | Description |
|------|---------|-------------|
| `--ref <branch-or-tag>` | default branch | Branch or tag to list commits from |
| `--limit <n>` | 20 | Max results |

### `gitlab commit view <project> <sha>`

View a commit's details.

```
orbit -p myprofile gl commit view 595 abc1234
```

---

## mr (merge request)

Alias: `merge-request`

### `gitlab mr list <project>`

List merge requests.

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | | Filter: `opened`, `closed`, `merged`, `all` |
| `--limit <n>` | 20 | Max results |

**Output columns:** IID, State (shows "draft" for draft MRs), Title, Author

### `gitlab mr view <project> <mr-iid>`

View merge request details.

**Output fields:** IID, Title, State, Source Branch, Target Branch, Author, Assignee, Labels, Merge Status, Has Conflicts, Comments Count, URL, Description

### `gitlab mr create <project>`

Create a merge request.

| Flag | Required | Description |
|------|----------|-------------|
| `--source <branch>` | Yes | Source branch |
| `--target <branch>` | Yes | Target branch |
| `--title <text>` | Yes | MR title |
| `--description <text>` | No | MR description |

```
orbit -p myprofile gl mr create 595 \
  --source feature/login --target main --title "Add login page"
```

### `gitlab mr close <project> <mr-iid>`

Close a merge request.

```
orbit -p myprofile gl mr close 595 42
```

### `gitlab mr merge <project> <mr-iid>`

Merge a merge request.

| Flag | Default | Description |
|------|---------|-------------|
| `--squash` | false | Squash commits |

### `gitlab mr comment <project> <mr-iid>`

Add a comment to a merge request.

| Flag | Required | Description |
|------|----------|-------------|
| `--body <text>` | Yes | Comment body |

### `gitlab mr notes <project> <mr-iid>`

List comments on a merge request (excludes system-generated notes).

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `gitlab mr approve <project> <mr-iid>`

Approve a merge request.

```
orbit -p myprofile gl mr approve 595 42
```

### `gitlab mr unapprove <project> <mr-iid>`

Revoke approval. Alias: `revoke`.

```
orbit -p myprofile gl mr unapprove 595 42
```

### `gitlab mr approvals <project> <mr-iid>`

View approval status — required, remaining, who approved.

```
orbit -p myprofile gl mr approvals 595 42
```

### `gitlab mr rebase <project> <mr-iid>`

Rebase onto the target branch.

```
orbit -p myprofile gl mr rebase 595 42
```

### `gitlab mr reopen <project> <mr-iid>`

Reopen a closed merge request.

```
orbit -p myprofile gl mr reopen 595 42
```

### `gitlab mr edit <project> <mr-iid>`

Edit a merge request.

| Flag | Description |
|------|-------------|
| `--title` | New title |
| `--description` | New description |
| `--assignee` | Assignee user ID |
| `--labels` | Comma-separated labels |
| `--target` | Target branch |
| `--draft` | Mark as draft |

```
orbit -p myprofile gl mr edit 595 42 --title "Updated title"
orbit -p myprofile gl mr edit 595 42 --assignee 15 --labels "bug,urgent"
```

---

## pipeline

Aliases: `pipe`, `ci`

### `gitlab pipeline list <project>`

List pipelines.

| Flag | Default | Description |
|------|---------|-------------|
| `--ref <branch-or-tag>` | | Filter by branch/tag |
| `--status <status>` | | Filter: `running`, `pending`, `success`, `failed`, `canceled` |
| `--limit <n>` | 20 | Max results |

### `gitlab pipeline view <project> <pipeline-id>`

View pipeline details.

**Output fields:** ID, Status, Ref, SHA, Source, Created, URL

### `gitlab pipeline jobs <project> <pipeline-id>`

List jobs in a pipeline.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

**Output columns:** ID, Name, Stage, Status, Duration

### `gitlab pipeline retry <project> <pipeline-id>`

Retry a pipeline (re-runs failed jobs).

### `gitlab pipeline cancel <project> <pipeline-id>`

Cancel a running pipeline.

---

## schedule

### `gitlab schedule list <project>`

List pipeline schedules.

```
orbit -p myprofile gl schedule list 595
```

### `gitlab schedule view <project> <schedule-id>`

View schedule details.

```
orbit -p myprofile gl schedule view 595 42
```

### `gitlab schedule create <project>`

Create a pipeline schedule.

| Flag | Default | Description |
|------|---------|-------------|
| `--desc` | | Schedule description (required) |
| `--ref` | `main` | Branch or tag to run |
| `--cron` | | Cron expression (required) |
| `--timezone` | `UTC` | Cron timezone |

```
orbit -p myprofile gl schedule create 595 --desc "Nightly build" --ref main --cron "0 2 * * *"
orbit -p myprofile gl schedule create 595 --desc "Weekly deploy" --cron "0 9 * * 1" --timezone "America/Mexico_City"
```

### `gitlab schedule update <project> <schedule-id>`

Update a pipeline schedule.

| Flag | Description |
|------|-------------|
| `--desc` | Schedule description |
| `--ref` | Branch or tag |
| `--cron` | Cron expression |
| `--timezone` | Cron timezone |
| `--active` | Enable/disable (`true`/`false`) |

```
orbit -p myprofile gl schedule update 595 42 --cron "0 3 * * *"
orbit -p myprofile gl schedule update 595 42 --active=false
```

### `gitlab schedule delete <project> <schedule-id>`

Delete a pipeline schedule.

```
orbit -p myprofile gl schedule delete 595 42
```

### `gitlab schedule run <project> <schedule-id>`

Trigger a scheduled pipeline immediately.

```
orbit -p myprofile gl schedule run 595 42
```

### `gitlab schedule var <project> <schedule-id> <key> <value>`

Add a variable to a schedule.

```
orbit -p myprofile gl schedule var 595 42 DEPLOY_ENV production
```

### `gitlab schedule var-delete <project> <schedule-id> <key>`

Delete a variable from a schedule.

```
orbit -p myprofile gl schedule var-delete 595 42 DEPLOY_ENV
```

---

## issue

### `gitlab issue list <project>`

List issues.

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | | Filter: `opened`, `closed`, `all` |
| `--labels <labels>` | | Filter by labels (comma-separated) |
| `--limit <n>` | 20 | Max results |

### `gitlab issue view <project> <issue-iid>`

View issue details.

**Output fields:** IID, Title, State, Author, Assignees, Labels, Due Date, Milestone, URL, Description

### `gitlab issue create <project>`

Create an issue.

| Flag | Required | Description |
|------|----------|-------------|
| `--title <text>` | Yes | Issue title |
| `--description <text>` | No | Issue description |
| `--labels <labels>` | No | Labels (comma-separated) |

### `gitlab issue close <project> <issue-iid>`

Close an issue.

---

## member

### `gitlab member list <project>`

List project members with access levels.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

**Access levels:** Guest (10), Reporter (20), Developer (30), Maintainer (40), Owner (50)

---

## user

### `gitlab user me`

Show the currently authenticated user.

**Output fields:** ID, Username, Name, Email, State, URL

### `gitlab user list`

List users.

| Flag | Default | Description |
|------|---------|-------------|
| `--search <text>` | | Search by username or name |
| `--limit <n>` | 20 | Max results |
