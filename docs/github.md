# Orbit GitHub Command Reference

Interact with GitHub repositories, pull requests, issues, actions, and more from the command line.

**Top-level command:** `orbit github` (alias: `gh`)

**Persistent flag (all subcommands):**

| Flag | Description |
|------|-------------|
| `--service` | GitHub service name, required only when a profile has multiple GitHub services configured |

**Notes:**

- The default API base URL is `https://api.github.com`.
- For GitHub Enterprise, set a custom base URL with `--base-url`.
- Repository arguments always use the `owner/repo` format (e.g., `octocat/hello-world`).

---

## Table of Contents

- [repo — View a repository](#repo)
- [repos — List repositories](#repos)
- [branch — Branch commands](#branch)
  - [branch list](#branch-list)
  - [branch view](#branch-view)
- [tag — Tag commands](#tag)
  - [tag list](#tag-list)
- [commit — Commit commands](#commit)
  - [commit list](#commit-list)
  - [commit view](#commit-view)
- [pr — Pull request commands](#pr)
  - [pr list](#pr-list)
  - [pr view](#pr-view)
  - [pr create](#pr-create)
  - [pr merge](#pr-merge)
  - [pr comment](#pr-comment)
  - [pr comments](#pr-comments)
- [issue — Issue commands](#issue)
  - [issue list](#issue-list)
  - [issue view](#issue-view)
  - [issue create](#issue-create)
  - [issue comment](#issue-comment)
  - [issue close](#issue-close)
- [release — Release commands](#release)
  - [release list](#release-list)
  - [release view](#release-view)
  - [release latest](#release-latest)
- [run — Workflow run commands](#run)
  - [run list](#run-list)
  - [run view](#run-view)
  - [run rerun](#run-rerun)
  - [run cancel](#run-cancel)
  - [run watch](#run-watch)
- [secret — Repository secret commands](#secret)
  - [secret list](#secret-list)
  - [secret set](#secret-set)
  - [secret delete](#secret-delete)
- [user — User commands](#user)
  - [user me](#user-me)
  - [user view](#user-view)

---

## repo

View details of a repository.

```
orbit github repo [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Repository to view |

**Example:**

```bash
orbit github repo octocat/hello-world -p myprofile
```

---

## repos

List repositories accessible to the authenticated user.

```
orbit github repos [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--search` | | Filter repositories by search query |
| `--limit` | 50 | Maximum number of repositories to return |

**Examples:**

```bash
# List repositories
orbit github repos -p myprofile

# Search for repositories matching a keyword
orbit github repos --search "api" -p myprofile

# Limit results
orbit github repos --limit 10 -p myprofile
```

---

## repo edit

Edit repository settings.

```
orbit github repo edit [owner/repo] [flags] -p myprofile
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--description` | string | New repository description. |
| `--private` | bool | Set private status. |
| `--default-branch` | string | Set default branch. |
| `--archived` | bool | Archive or unarchive the repository. |

**Examples:**

```bash
# Archive a repository
orbit gh repo edit Paybook/ai --archived -p myprofile

# Unarchive a repository
orbit gh repo edit Paybook/ai --archived=false -p myprofile

# Update description
orbit gh repo edit Paybook/ai --description "Updated description" -p myprofile
```

---

## repo collaborator

Manage repository collaborators (alias: `collab`).

### repo collaborator list

List repository collaborators with permissions.

```bash
orbit gh repo collab list Paybook/ai -p myprofile
```

### repo collaborator add

Add a collaborator to a repository.

```
orbit github repo collaborator add [owner/repo] [username] [flags] -p myprofile
```

| Flag | Default | Description |
|------|---------|-------------|
| `--permission` | `push` | Permission: `pull`, `triage`, `push`, `maintain`, `admin`. |

```bash
orbit gh repo collab add Paybook/ai jorgemuza --permission admin -p myprofile
```

### repo collaborator remove

Remove a collaborator from a repository.

```bash
orbit gh repo collab remove Paybook/ai jorgemuza -p myprofile
```

---

## branch

Manage repository branches.

### branch list

List branches in a repository.

**Aliases:** `ls`

```
orbit github branch list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of branches to return |

**Examples:**

```bash
orbit github branch list octocat/hello-world -p myprofile
orbit github branch ls octocat/hello-world --limit 20 -p myprofile
```

### branch view

View details of a specific branch.

```
orbit github branch view [owner/repo] [branch] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `branch` | Branch name to view |

**Example:**

```bash
orbit github branch view octocat/hello-world main -p myprofile
```

---

## tag

Manage repository tags.

### tag list

List tags in a repository.

**Aliases:** `ls`

```
orbit github tag list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of tags to return |

**Example:**

```bash
orbit github tag list octocat/hello-world --limit 10 -p myprofile
```

---

## commit

Manage repository commits.

### commit list

List commits in a repository.

**Aliases:** `ls`

```
orbit github commit list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of commits to return |

**Example:**

```bash
orbit github commit list octocat/hello-world --limit 20 -p myprofile
```

### commit view

View details of a specific commit.

```
orbit github commit view [owner/repo] [sha] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `sha` | Full or abbreviated commit SHA |

**Example:**

```bash
orbit github commit view octocat/hello-world abc1234 -p myprofile
```

---

## pr

Manage pull requests.

**Aliases:** `pull-request`

### pr list

List pull requests in a repository.

**Aliases:** `ls`

```
orbit github pr list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | | Filter by state: `open`, `closed`, or `all` |
| `--limit` | | Maximum number of pull requests to return |

**Examples:**

```bash
orbit github pr list octocat/hello-world -p myprofile
orbit github pr ls octocat/hello-world --state closed --limit 5 -p myprofile
```

### pr view

View details of a specific pull request.

```
orbit github pr view [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

**Example:**

```bash
orbit github pr view octocat/hello-world 42 -p myprofile
```

### pr create

Create a new pull request.

```
orbit github pr create [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Required | Description |
|------|----------|-------------|
| `--head` | Yes | Source branch name |
| `--base` | Yes | Target branch name |
| `--title` | Yes | Pull request title |
| `--body` | No | Pull request description |

**Examples:**

```bash
# Minimal create
orbit github pr create octocat/hello-world \
  --head feature/login \
  --base main \
  --title "Add login page" \
  -p myprofile

# With body
orbit github pr create octocat/hello-world \
  --head feature/login \
  --base main \
  --title "Add login page" \
  --body "Implements the new login flow with OAuth support." \
  -p myprofile
```

### pr merge

Merge a pull request.

```
orbit github pr merge [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Default | Description |
|------|---------|-------------|
| `--squash` | | Squash commits before merging |

**Examples:**

```bash
orbit github pr merge octocat/hello-world 42 -p myprofile
orbit github pr merge octocat/hello-world 42 --squash -p myprofile
```

### pr comment

Add a comment to a pull request.

```
orbit github pr comment [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Required | Description |
|------|----------|-------------|
| `--body` | Yes | Comment text |

**Example:**

```bash
orbit github pr comment octocat/hello-world 42 \
  --body "LGTM, approving." \
  -p myprofile
```

### pr comments

List comments on a pull request.

```
orbit github pr comments [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of comments to return |

**Example:**

```bash
orbit github pr comments octocat/hello-world 42 --limit 10 -p myprofile
```

---

## issue

Manage repository issues.

### issue list

List issues in a repository.

**Aliases:** `ls`

```
orbit github issue list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | | Filter by state: `open`, `closed`, or `all` |
| `--labels` | | Filter by comma-separated label names |
| `--limit` | | Maximum number of issues to return |

**Examples:**

```bash
orbit github issue list octocat/hello-world -p myprofile
orbit github issue ls octocat/hello-world --state open --labels "bug,critical" --limit 25 -p myprofile
```

### issue view

View details of a specific issue.

```
orbit github issue view [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

**Example:**

```bash
orbit github issue view octocat/hello-world 7 -p myprofile
```

### issue create

Create a new issue.

```
orbit github issue create [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Issue title |
| `--body` | No | Issue description |
| `--labels` | No | Comma-separated label names |

**Examples:**

```bash
orbit github issue create octocat/hello-world \
  --title "Fix broken login" \
  -p myprofile

orbit github issue create octocat/hello-world \
  --title "Fix broken login" \
  --body "The login form returns a 500 error when submitting valid credentials." \
  --labels "bug,high-priority" \
  -p myprofile
```

### issue comment

Add a comment to an issue.

```
orbit github issue comment [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

| Flag | Required | Description |
|------|----------|-------------|
| `--body` | Yes | Comment text |

**Example:**

```bash
orbit github issue comment octocat/hello-world 7 \
  --body "Reproduced on latest main. Investigating." \
  -p myprofile
```

### issue close

Close an issue.

```
orbit github issue close [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

**Example:**

```bash
orbit github issue close octocat/hello-world 7 -p myprofile
```

---

## release

Manage repository releases.

### release list

List releases in a repository.

**Aliases:** `ls`

```
orbit github release list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of releases to return |

**Example:**

```bash
orbit github release list octocat/hello-world --limit 5 -p myprofile
```

### release view

View details of a specific release by tag.

```
orbit github release view [owner/repo] [tag] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `tag` | Release tag name |

**Example:**

```bash
orbit github release view octocat/hello-world v1.2.0 -p myprofile
```

### release latest

View the latest release of a repository.

```
orbit github release latest [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

**Example:**

```bash
orbit github release latest octocat/hello-world -p myprofile
```

---

## run

Manage GitHub Actions workflow runs.

**Aliases:** `actions`

### run list

List workflow runs in a repository.

**Aliases:** `ls`

```
orbit github run list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of runs to return |

**Example:**

```bash
orbit github run list octocat/hello-world --limit 10 -p myprofile
```

### run view

View details of a specific workflow run.

```
orbit github run view [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run view octocat/hello-world 123456789 -p myprofile
```

### run rerun

Re-run a workflow run.

```
orbit github run rerun [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run rerun octocat/hello-world 123456789 -p myprofile
```

### run cancel

Cancel an in-progress workflow run.

```
orbit github run cancel [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run cancel octocat/hello-world 123456789 -p myprofile
```

### run watch

Watch a workflow run until it completes.

```
orbit github run watch [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run watch octocat/hello-world 123456789 -p myprofile
```

---

## secret

Manage repository secrets.

### secret list

List secrets configured on a repository.

**Aliases:** `ls`

```
orbit github secret list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

**Example:**

```bash
orbit github secret list octocat/hello-world -p myprofile
```

### secret set

Create or update a repository secret.

```
orbit github secret set [owner/repo] [secret-name] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `secret-name` | Name of the secret |

| Flag | Required | Description |
|------|----------|-------------|
| `--value` | Yes | Secret value |

**Example:**

```bash
orbit github secret set octocat/hello-world API_KEY \
  --value "sk-abc123" \
  -p myprofile
```

### secret delete

Delete a repository secret.

```
orbit github secret delete [owner/repo] [secret-name] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `secret-name` | Name of the secret to delete |

**Example:**

```bash
orbit github secret delete octocat/hello-world API_KEY -p myprofile
```

---

## user

View GitHub user information.

### user me

Display the currently authenticated user.

```
orbit github user me [flags]
```

**Example:**

```bash
orbit github user me -p myprofile
```

### user view

View a GitHub user's profile.

```
orbit github user view [username] [flags]
```

| Argument | Description |
|----------|-------------|
| `username` | GitHub username to look up |

**Example:**

```bash
orbit github user view octocat -p myprofile
```
