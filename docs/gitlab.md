# Orbit GitLab Command Reference

Complete reference for the `orbit gitlab` (alias: `gl`) commands.

Default base URL: `https://gitlab.com`. For self-hosted instances, use `--base-url` in your profile configuration.

Projects can be specified by numeric ID or full path (e.g., `schools/frontend/my-app`).

---

## Table of Contents

- [Global Flags](#global-flags)
- [project](#project)
  - [project view](#project-view)
  - [project create](#project-create)
- [projects](#projects)
- [group](#group)
- [branch](#branch)
- [tag](#tag)
- [commit](#commit)
- [mr (merge-request)](#mr)
- [pipeline](#pipeline)
- [job](#job)
- [runner](#runner)
- [issue](#issue)
- [file](#file)
- [member](#member)
- [user](#user)
- [variable](#variable)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--service` | GitLab service name. Use when a profile has multiple GitLab services configured. |

All examples below use `-p myprofile` to specify the Orbit profile.

---

## project

Manage projects.

### project view

View details of a single project.

```
orbit gitlab project view [id-or-path] -p myprofile
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `id-or-path` | Project ID (numeric) or full path. |

**Examples:**

```bash
# View project by ID
orbit gl project view 42 -p myprofile

# View project by path
orbit gl project view schools/frontend/my-app -p myprofile
```

### project create

Create a new project.

```
orbit gitlab project create [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--name` | string | | Project name. | Yes |
| `--path` | string | | Project path (defaults to slugified name). | No |
| `--visibility` | string | `private` | Visibility: `private`, `internal`, `public`. | No |
| `--description` | string | | Project description. | No |
| `--namespace` | string | | Namespace ID or group path (e.g. `foundation`). | No |

**Examples:**

```bash
# Create a private project
orbit gl project create --name "my-app" -p myprofile

# Create a project under a group namespace
orbit gl project create --name "my-app" --namespace foundation -p myprofile

# Create a project with namespace ID and public visibility
orbit gl project create --name "my-app" --namespace 42 --visibility public -p myprofile

# Create a project with a description
orbit gl project create --name "my-app" --description "My new project" -p myprofile
```

---

## project list

List projects.

```
orbit gitlab project list [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--search` | string | | Search projects by name. |
| `--group` | string | | List projects in a group (ID or path). |
| `--limit` | int | 50 | Maximum number of results. |

**Examples:**

```bash
# List all projects
orbit gl project list -p myprofile

# Search for projects by name
orbit gl project list --search "frontend" -p myprofile

# List projects in a group
orbit gl project list --group "schools/frontend" -p myprofile

# Limit results
orbit gl project list --search "api" --limit 10 -p myprofile
```

---

## project activity

List projects with recent activity, showing which branches received pushes.

```
orbit gitlab project activity [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--since` | string | | Show activity since date (YYYY-MM-DD). |
| `--days` | int | | Show activity from the last N days. |
| `--group` | string | | Filter to projects in a group (ID or path). |
| `--branch` | string | | Filter by branch names (comma-separated). |
| `--limit` | int | 100 | Maximum projects to scan. |

When neither `--since` nor `--days` is provided, defaults to today.

**Examples:**

```bash
# Projects with activity today
orbit gl project activity -p myprofile

# Activity in the last 7 days
orbit gl project activity --days 7 -p myprofile

# Activity since a specific date
orbit gl project activity --since 2026-03-20 -p myprofile

# Scoped to a group
orbit gl project activity --group schools --days 3 -p myprofile

# Only show activity on specific branches
orbit gl project activity --branch development,master -p myprofile
```

---

## group

Manage and view GitLab groups.

### group list

List all groups. Alias: `ls`.

```bash
orbit gl group list -p myprofile
orbit gl group ls -p myprofile
```

### group view

View details of a group.

```
orbit gitlab group view [id-or-path] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `id-or-path` | Group ID (numeric) or full path. |

```bash
orbit gl group view schools -p myprofile
orbit gl group view 15 -p myprofile
```

### group subgroups

List subgroups of a group.

```
orbit gitlab group subgroups [id-or-path] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `id-or-path` | Group ID (numeric) or full path. |

```bash
orbit gl group subgroups schools -p myprofile
```

---

## branch

Manage branches in a project.

### branch list

List branches. Alias: `ls`.

```
orbit gitlab branch list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl branch ls schools/frontend/my-app -p myprofile
orbit gl branch list 42 --limit 10 -p myprofile
```

### branch view

View details of a specific branch.

```
orbit gitlab branch view [project] [branch] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `branch` | Branch name. |

```bash
orbit gl branch view schools/frontend/my-app feature/login -p myprofile
```

### branch create

Create a new branch.

```
orbit gitlab branch create [project] [branch] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `branch` | Name for the new branch. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--ref` | string | | Source ref (branch, tag, or SHA) to branch from. |

```bash
# Create a branch from main
orbit gl branch create schools/frontend/my-app feature/signup --ref main -p myprofile

# Create a branch from a tag
orbit gl branch create 42 hotfix/payment --ref v1.2.0 -p myprofile
```

### branch delete

Delete a branch.

```
orbit gitlab branch delete [project] [branch] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `branch` | Branch name to delete. |

```bash
orbit gl branch delete schools/frontend/my-app feature/old-feature -p myprofile
```

---

## tag

Manage tags in a project.

### tag list

List tags. Alias: `ls`.

```
orbit gitlab tag list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl tag ls schools/frontend/my-app -p myprofile
orbit gl tag list 42 --limit 5 -p myprofile
```

### tag create

Create a new tag.

```
orbit gitlab tag create [project] [tag] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `tag` | Tag name. |

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--ref` | string | | Source ref to tag (branch, SHA). | Yes |
| `--message` | string | | Tag message (annotated tag). | No |

```bash
# Create a lightweight tag
orbit gl tag create schools/frontend/my-app v2.0.0 --ref main -p myprofile

# Create an annotated tag with a message
orbit gl tag create 42 v2.0.0 --ref main --message "Release 2.0.0" -p myprofile
```

---

## commit

View commits in a project.

### commit list

List commits. Alias: `ls`.

```
orbit gitlab commit list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl commit ls schools/frontend/my-app -p myprofile
orbit gl commit list 42 --limit 20 -p myprofile
```

### commit view

View details of a specific commit.

```
orbit gitlab commit view [project] [sha] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `sha` | Commit SHA. |

```bash
orbit gl commit view schools/frontend/my-app abc123def -p myprofile
```

---

## mr

Manage merge requests. Aliases: `merge-request`.

### mr list

List merge requests. Alias: `ls`.

```
orbit gitlab mr list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | | Filter by state: `opened`, `closed`, `merged`, `all`. |
| `--limit` | int | 20 | Maximum number of results. |

```bash
# List open merge requests
orbit gl mr ls schools/frontend/my-app --state opened -p myprofile

# List all merge requests
orbit gl mr list 42 --state all --limit 50 -p myprofile
```

### mr view

View details of a merge request.

```
orbit gitlab mr view [project] [mr-iid] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `mr-iid` | Merge request IID (project-scoped ID). |

```bash
orbit gl mr view schools/frontend/my-app 15 -p myprofile
```

### mr create

Create a new merge request.

```
orbit gitlab mr create [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--source` | string | | Source branch. | Yes |
| `--target` | string | | Target branch. | Yes |
| `--title` | string | | Merge request title. | Yes |
| `--description` | string | | Merge request description. | No |

```bash
orbit gl mr create schools/frontend/my-app \
  --source feature/login \
  --target main \
  --title "Add login page" \
  --description "Implements the login flow with OAuth support" \
  -p myprofile
```

### mr merge

Merge a merge request.

```
orbit gitlab mr merge [project] [mr-iid] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `mr-iid` | Merge request IID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--squash` | bool | | Squash commits on merge. |

```bash
orbit gl mr merge schools/frontend/my-app 15 -p myprofile
orbit gl mr merge 42 15 --squash -p myprofile
```

### mr comment

Add a comment to a merge request.

```
orbit gitlab mr comment [project] [mr-iid] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `mr-iid` | Merge request IID. |

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--body` | string | | Comment text. | Yes |

```bash
orbit gl mr comment schools/frontend/my-app 15 --body "LGTM, approved!" -p myprofile
```

### mr notes

List notes (comments) on a merge request.

```
orbit gitlab mr notes [project] [mr-iid] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `mr-iid` | Merge request IID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Maximum number of results. |

```bash
orbit gl mr notes schools/frontend/my-app 15 -p myprofile
orbit gl mr notes 42 15 --limit 10 -p myprofile
```

---

## pipeline

Manage CI/CD pipelines. Aliases: `pipe`, `ci`.

### pipeline list

List pipelines. Alias: `ls`.

```
orbit gitlab pipeline list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl ci ls schools/frontend/my-app -p myprofile
orbit gl pipe list 42 --limit 5 -p myprofile
```

### pipeline view

View details of a pipeline.

```
orbit gitlab pipeline view [project] [pipeline-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `pipeline-id` | Pipeline ID. |

```bash
orbit gl ci view schools/frontend/my-app 98765 -p myprofile
```

### pipeline jobs

List jobs in a pipeline.

```
orbit gitlab pipeline jobs [project] [pipeline-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `pipeline-id` | Pipeline ID. |

```bash
orbit gl ci jobs schools/frontend/my-app 98765 -p myprofile
```

### pipeline retry

Retry a pipeline.

```
orbit gitlab pipeline retry [project] [pipeline-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `pipeline-id` | Pipeline ID. |

```bash
orbit gl ci retry schools/frontend/my-app 98765 -p myprofile
```

### pipeline cancel

Cancel a running pipeline.

```
orbit gitlab pipeline cancel [project] [pipeline-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `pipeline-id` | Pipeline ID. |

```bash
orbit gl ci cancel schools/frontend/my-app 98765 -p myprofile
```

---

## job

Manage individual CI/CD jobs.

### job list

List jobs. Alias: `ls`.

```
orbit gitlab job list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl job ls schools/frontend/my-app -p myprofile
```

### job view

View details of a job.

```
orbit gitlab job view [project] [job-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `job-id` | Job ID. |

```bash
orbit gl job view schools/frontend/my-app 123456 -p myprofile
```

### job log

View the log output of a job.

```
orbit gitlab job log [project] [job-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `job-id` | Job ID. |

```bash
orbit gl job log schools/frontend/my-app 123456 -p myprofile
```

### job play

Trigger a manual job.

```
orbit gitlab job play [project] [job-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `job-id` | Job ID. |

```bash
orbit gl job play schools/frontend/my-app 123456 -p myprofile
```

### job retry

Retry a failed job.

```
orbit gitlab job retry [project] [job-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `job-id` | Job ID. |

```bash
orbit gl job retry schools/frontend/my-app 123456 -p myprofile
```

### job cancel

Cancel a running job.

```
orbit gitlab job cancel [project] [job-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `job-id` | Job ID. |

```bash
orbit gl job cancel schools/frontend/my-app 123456 -p myprofile
```

---

## runner

Manage GitLab runners.

### runner list

List runners for a project. Alias: `ls`.

```
orbit gitlab runner list [project] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

```bash
orbit gl runner ls schools/frontend/my-app -p myprofile
```

### runner list-all

List all runners available to the instance.

```
orbit gitlab runner list-all -p myprofile
```

```bash
orbit gl runner list-all -p myprofile
```

### runner view

View details of a runner.

```
orbit gitlab runner view [runner-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `runner-id` | Runner ID. |

```bash
orbit gl runner view 7 -p myprofile
```

### runner enable

Enable a runner.

```
orbit gitlab runner enable [runner-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `runner-id` | Runner ID. |

```bash
orbit gl runner enable 7 -p myprofile
```

### runner disable

Disable a runner.

```
orbit gitlab runner disable [runner-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `runner-id` | Runner ID. |

```bash
orbit gl runner disable 7 -p myprofile
```

---

## issue

Manage project issues.

### issue list

List issues. Alias: `ls`.

```
orbit gitlab issue list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | | Filter by state: `opened`, `closed`, `all`. |
| `--labels` | string | | Filter by labels (comma-separated). |
| `--limit` | int | 20 | Maximum number of results. |

```bash
# List open issues
orbit gl issue ls schools/frontend/my-app --state opened -p myprofile

# Filter by labels
orbit gl issue list 42 --labels "bug,critical" --limit 10 -p myprofile
```

### issue view

View details of an issue.

```
orbit gitlab issue view [project] [issue-iid] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `issue-iid` | Issue IID (project-scoped ID). |

```bash
orbit gl issue view schools/frontend/my-app 8 -p myprofile
```

### issue create

Create a new issue.

```
orbit gitlab issue create [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--title` | string | | Issue title. | Yes |
| `--description` | string | | Issue description. | No |
| `--labels` | string | | Comma-separated labels. | No |

```bash
orbit gl issue create schools/frontend/my-app \
  --title "Fix broken navigation on mobile" \
  --description "The hamburger menu does not open on iOS Safari" \
  --labels "bug,mobile" \
  -p myprofile
```

### issue close

Close an issue.

```
orbit gitlab issue close [project] [issue-iid] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `issue-iid` | Issue IID. |

```bash
orbit gl issue close schools/frontend/my-app 8 -p myprofile
```

---

## file

Read and update files in a repository.

### file read

Read the contents of a file.

```
orbit gitlab file read [project] [file-path] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `file-path` | Path to the file in the repository. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--ref` | string | main | Branch, tag, or commit SHA to read from. |

```bash
# Read from default branch (main)
orbit gl file read schools/frontend/my-app src/index.ts -p myprofile

# Read from a specific branch
orbit gl file read 42 README.md --ref develop -p myprofile
```

### file update

Update a file in the repository.

```
orbit gitlab file update [project] [file-path] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `file-path` | Path to the file in the repository. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--message` | string | | Commit message. |
| `--content` | string | | New file content. |
| `--ref` | string | main | Branch to commit to. |

```bash
orbit gl file update schools/frontend/my-app config.json \
  --content '{"version": "2.0"}' \
  --message "Update config to v2" \
  --ref main \
  -p myprofile
```

---

## member

View project members.

### member list

List members of a project. Alias: `ls`.

```
orbit gitlab member list [project] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl member ls schools/frontend/my-app -p myprofile
```

---

## user

Manage and view GitLab users.

### user list

List users. Alias: `ls`.

```
orbit gitlab user list [flags] -p myprofile
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--search` | string | | Search users by name or username. |
| `--limit` | int | | Maximum number of results. |

```bash
orbit gl user ls --search "jorge" -p myprofile
orbit gl user list --limit 10 -p myprofile
```

### user me

View the currently authenticated user.

```
orbit gitlab user me -p myprofile
```

```bash
orbit gl user me -p myprofile
```

### user view

View details of a specific user.

```
orbit gitlab user view [user-id] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `user-id` | User ID. |

```bash
orbit gl user view 12 -p myprofile
```

---

## variable

Manage CI/CD variables for a project.

### variable list

List all variables. Alias: `ls`.

```
orbit gitlab variable list [project] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |

```bash
orbit gl variable ls schools/frontend/my-app -p myprofile
```

### variable get

Get the value of a specific variable.

```
orbit gitlab variable get [project] [variable-key] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `variable-key` | Variable key name. |

```bash
orbit gl variable get schools/frontend/my-app DATABASE_URL -p myprofile
```

### variable set

Set (create or update) a variable.

```
orbit gitlab variable set [project] [variable-key] [flags] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `variable-key` | Variable key name. |

| Flag | Type | Default | Description | Required |
|------|------|---------|-------------|----------|
| `--value` | string | | Variable value. | Yes |

```bash
orbit gl variable set schools/frontend/my-app API_KEY --value "sk-abc123" -p myprofile
```

### variable delete

Delete a variable.

```
orbit gitlab variable delete [project] [variable-key] -p myprofile
```

| Argument | Description |
|----------|-------------|
| `project` | Project ID or path. |
| `variable-key` | Variable key name. |

```bash
orbit gl variable delete schools/frontend/my-app OLD_SECRET -p myprofile
```
