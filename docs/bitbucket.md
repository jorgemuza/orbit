# Orbit Bitbucket Command Reference

Manage Bitbucket resources from the command line.

**Command:** `orbit bitbucket`
**Alias:** `orbit bb`

## Table of Contents

- [Global Flags](#global-flags)
- [API Variants](#api-variants)
- [Projects](#project)
- [Repositories](#repo)
- [Branches](#branch)
- [Tags](#tag)
- [Commits](#commit)
- [Pull Requests](#pr)
- [Users](#user)
- [Reviewer Conditions](#reviewer-condition)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--service` | Bitbucket service name. Required only when a profile has multiple Bitbucket services configured. |
| `--debug` | Print HTTP request/response details to stderr for troubleshooting. |

All examples below use `-p myprofile` to specify the Orbit profile.

---

## API Variants

Orbit supports two Bitbucket backends:

| Variant | API | Base URL |
|---------|-----|----------|
| **Cloud** | Bitbucket Cloud REST API v2.0 | `https://api.bitbucket.org/2.0` (default) |
| **Server** | Bitbucket Data Center REST API | Requires `--base-url` |

- **Cloud** uses `workspace/repo` format for addressing repositories.
- **Server** uses `project-key/repo-slug` format.

---

## project

Manage Bitbucket projects.

### `project list`

List projects.

**Alias:** `project ls`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of projects to return. |

```bash
# List all projects
orbit bb project list -p myprofile

# List with a limit
orbit bb project list -p myprofile --limit 10
```

### `project view`

View details of a single project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key (e.g. `TEAM`). |

```bash
orbit bb project view TEAM -p myprofile
```

---

## repo

Manage repositories within a project.

### `repo list`

List repositories in a project.

**Alias:** `repo ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of repositories to return. |

```bash
# List all repos in a project
orbit bb repo list TEAM -p myprofile

# Limit results
orbit bb repo list TEAM -p myprofile --limit 5
```

### `repo view`

View details of a single repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

```bash
orbit bb repo view TEAM my-service -p myprofile
```

---

## branch

Manage branches in a repository.

### `branch list`

List branches.

**Alias:** `branch ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of branches to return. |

```bash
orbit bb branch list TEAM my-service -p myprofile
orbit bb branch ls TEAM my-service -p myprofile --limit 20
```

### `branch create`

Create a new branch.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `branch` | 3 | Name of the new branch. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--ref` | string | Yes | The commit SHA or branch name to branch from. |

```bash
# Create a feature branch from main
orbit bb branch create TEAM my-service feature/login --ref main -p myprofile

# Create a branch from a specific commit
orbit bb branch create TEAM my-service hotfix/urgent --ref a1b2c3d -p myprofile
```

### `branch delete`

Delete a branch.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `branch` | 3 | Name of the branch to delete. |

```bash
orbit bb branch delete TEAM my-service feature/login -p myprofile
```

### `branch default`

View or set the default branch of a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--set` | string | No | Set the default branch to this value. Omit to view the current default. |

```bash
# View the current default branch
orbit bb branch default TEAM my-service -p myprofile

# Change the default branch to develop
orbit bb branch default TEAM my-service -p myprofile --set develop
```

---

## tag

Manage tags in a repository.

### `tag list`

List tags.

**Alias:** `tag ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of tags to return. |

```bash
orbit bb tag list TEAM my-service -p myprofile
orbit bb tag ls TEAM my-service -p myprofile --limit 10
```

### `tag create`

Create a new tag.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `tag` | 3 | Name of the tag. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--ref` | string | Yes | The commit SHA or branch name to tag. |

```bash
# Tag the current HEAD of main
orbit bb tag create TEAM my-service v1.0.0 --ref main -p myprofile

# Tag a specific commit
orbit bb tag create TEAM my-service v1.0.1 --ref a1b2c3d -p myprofile
```

---

## commit

Browse commit history.

### `commit list`

List commits in a repository.

**Alias:** `commit ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of commits to return. |

```bash
orbit bb commit list TEAM my-service -p myprofile
orbit bb commit ls TEAM my-service -p myprofile --limit 50
```

### `commit view`

View details of a single commit.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `sha` | 3 | The commit SHA. |

```bash
orbit bb commit view TEAM my-service a1b2c3d4e5f6 -p myprofile
```

---

## pr

Manage pull requests.

### `pr list`

List pull requests in a repository.

**Alias:** `pr ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | — | Filter by state: `open`, `declined`, `merged`, or `all`. |
| `--limit` | int | — | Maximum number of pull requests to return. |

```bash
# List open PRs (default)
orbit bb pr list TEAM my-service -p myprofile

# List merged PRs
orbit bb pr list TEAM my-service -p myprofile --state merged

# List all PRs with a limit
orbit bb pr ls TEAM my-service -p myprofile --state all --limit 25
```

### `pr view`

View details of a single pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr view TEAM my-service 42 -p myprofile
```

### `pr create`

Create a new pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--source` | string | Yes | Source branch name. |
| `--target` | string | Yes | Target branch name. |
| `--title` | string | Yes | Pull request title. |
| `--description` | string | No | Pull request description body. |

```bash
# Minimal PR
orbit bb pr create TEAM my-service \
  --source feature/login \
  --target main \
  --title "Add login flow" \
  -p myprofile

# PR with description
orbit bb pr create TEAM my-service \
  --source feature/login \
  --target main \
  --title "Add login flow" \
  --description "Implements OAuth2 login with refresh token support." \
  -p myprofile
```

### `pr merge`

Merge a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr merge TEAM my-service 42 -p myprofile
```

### `pr decline`

Decline a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr decline TEAM my-service 42 -p myprofile
```

### `pr comment`

Add a comment to a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--body` | string | Yes | The comment text. |

```bash
orbit bb pr comment TEAM my-service 42 --body "LGTM, approved." -p myprofile
```

### `pr diff`

View the unified diff for a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--context` | int | 3 | Number of context lines around changes. |

```bash
orbit bb pr diff TEAM my-service 42 -p myprofile
orbit bb pr diff TEAM my-service 42 --context 10 -p myprofile
```

### `pr approve`

Approve a pull request as the current user.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr approve TEAM my-service 42 -p myprofile
```

### `pr unapprove`

Remove your approval from a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr unapprove TEAM my-service 42 -p myprofile
```

### `pr activity`

View the activity feed of a pull request (comments, approvals, status changes).

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr activity TEAM my-service 42 -p myprofile
```

Nested comment replies are shown with indentation — all reply threads are traversed recursively.

---

## user

Manage Bitbucket users.

### `user list`

List users.

**Alias:** `user ls`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of users to return. |

```bash
orbit bb user list -p myprofile
orbit bb user ls -p myprofile --limit 100
```

---

## reviewer-condition

Manage project-level default reviewer conditions. These auto-assign reviewers to PRs and enforce required approvals.

**Alias:** `rc`

### `reviewer-condition list`

List all default reviewer conditions for a project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

```bash
orbit bb reviewer-condition list TEAM -p myprofile
orbit bb rc list TEAM -p myprofile -o json
```

### `reviewer-condition update`

Update the required approvals count for a condition.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `condition-id` | 2 | The condition ID. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--required-approvals` | int | Yes | Number of required approvals (0 = no requirement). |

```bash
# Temporarily bypass required approvals
orbit bb rc update TEAM 1063 --required-approvals 0 -p myprofile

# Restore required approvals
orbit bb rc update TEAM 1063 --required-approvals 2 -p myprofile
```

### `reviewer-condition delete`

Delete a default reviewer condition permanently.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `condition-id` | 2 | The condition ID. |

```bash
orbit bb rc delete TEAM 1063 -p myprofile
```
