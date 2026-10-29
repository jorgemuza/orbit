# Bitbucket CLI Command Reference

Complete reference for all `orbit bitbucket` (alias: `bb`) commands and flags.

## Table of Contents

- [project](#project)
- [repo](#repo)
- [branch](#branch)
- [tag](#tag)
- [commit](#commit)
- [pr (pull request)](#pr-pull-request)
- [code (search)](#code)
- [user](#user)
- [reviewer-condition (admin)](#reviewer-condition-admin)

## Global Flags

These flags apply to all bitbucket commands:

| Flag | Description |
|------|-------------|
| `-p, --profile <name>` | Profile to use (required) |
| `-o, --output <format>` | Output format: `table` (default), `json`, `yaml` |
| `--service <name>` | Bitbucket service name (if profile has multiple) |

---

## project

### `bitbucket project view <project-key>`

View a single project.

```
orbit -p myprofile bb project view L3SUP
```

**Output fields:** Key, Name, Description, Public, Type

### `bitbucket project search <text>`

Search projects by display **name** (not key). **Alias:** `find`. **Server only.**

| Flag | Default | Description |
|------|---------|-------------|
| `--permission <p>` | | `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN` |
| `--limit <n>` | 25 | Max results |

```
orbit -p myprofile bb project search agentic
orbit -p myprofile bb project search "" --permission PROJECT_ADMIN
```

A project keyed `EPCAP` but named "Epsilon PeopleCloud Agentic Platform" will NOT match `EPCAP`. With a key in hand, use `project view EPCAP`.

### `bitbucket project list`

List all projects.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

---

## repo

### `bitbucket repo view <project-key> <repo-slug>`

View repository details including clone URLs.

```
orbit -p myprofile bb repo view L3SUP agents-sre
```

**Output fields:** ID, Slug, Name, Description, State, SCM, Forkable, Project, Clone URLs

### `bitbucket repo list <project-key>`

List repositories in a project.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `bitbucket repo search [text]`

Search repositories by name **across every project**. **Alias:** `find`. **Server only.** Use this instead of listing each project's repos and filtering locally.

| Flag | Default | Description |
|------|---------|-------------|
| `--project <key>` | | One project, by key |
| `--project-name <text>` | | Projects whose display name matches (exclusive with `--project`) |
| `--permission <p>` | | `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN` |
| `--visibility <v>` | | `public` or `private` |
| `--archived <s>` | | `ACTIVE`, `ARCHIVED`, `ALL` |
| `--limit <n>` | 25 | Max results |

```
orbit -p myprofile bb repo search agentic
orbit -p myprofile bb repo search agentic --project EPCAP
orbit -p myprofile bb repo search payments --permission REPO_WRITE
orbit -p myprofile bb repo search "" --project-name "Ad Server" --limit 100
```

`--project` takes a KEY. Bitbucket's own filter is `projectname` and matches the display name, so the raw API returns an empty page for a key; orbit resolves the key first and then matches on it. Paging follows the server's `nextPageStart`, never `start+limit`, because Bitbucket's permission filtering makes that cursor jump ahead of the page it returned.

---

## branch

### `bitbucket branch list <project-key> <repo-slug>`

List branches.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter branches by name |
| `--limit <n>` | 50 | Max results |

### `bitbucket branch default <project-key> <repo-slug>`

Show the default branch.

### `bitbucket branch create <project-key> <repo-slug> <name> <start-point>`

Create a branch from a ref (branch name, tag, or commit SHA).

```
orbit -p myprofile bb branch create L3SUP agents-sre feature/new-thing main
```

### `bitbucket branch delete <project-key> <repo-slug> <branch-name>`

Delete a branch. **Server only.** Either `feature/x` or `refs/heads/feature/x` works.

```
orbit -p myprofile bb branch delete L3SUP agents-sre feature/old-thing
```

Prints the commit the branch was at, which is how to recreate it: `Deleted branch feature/old-thing (was at 5e85a5908712)`.

Deletion lives in the `branch-utils` plugin, not the core API whose `/branches` is GET-only and answers a DELETE with 405 - orbit sent it to the wrong one until v0.72.0. `branch-utils` also answers a delete of a branch that does not exist with the same 204 it gives a real one, so orbit checks the branch exists first and checks it is gone afterwards; a name that is not there is an error, not a success line.

---

## tag

### `bitbucket tag list <project-key> <repo-slug>`

List tags.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter tags by name |
| `--limit <n>` | 50 | Max results |

### `bitbucket tag create <project-key> <repo-slug> <tag-name> <start-point>`

Create a tag from a ref.

| Flag | Description |
|------|-------------|
| `-m, --message <text>` | Tag message (annotated tag) |

```
orbit -p myprofile bb tag create L3SUP agents-sre v1.0.0 main -m "Release v1.0.0"
```

---

## commit

### `bitbucket commit list <project-key> <repo-slug>`

List recent commits.

| Flag | Default | Description |
|------|---------|-------------|
| `--branch <name>` | default branch | Branch to list commits from |
| `--limit <n>` | 20 | Max results |

### `bitbucket commit view <project-key> <repo-slug> <commit-id>`

View a commit's details.

```
orbit -p myprofile bb commit view L3SUP agents-sre abc1234def5678
```

**Output fields:** ID, Author, Date, Message

---

## pr (pull request)

Alias: `pull-request`

### `bitbucket pr list <project-key> <repo-slug>`

List pull requests.

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | | Filter: `OPEN`, `MERGED`, `DECLINED`, `ALL` |
| `--limit <n>` | 25 | Max results |

**Output columns:** ID, State, Author, Title

### `bitbucket pr view <project-key> <repo-slug> <pr-id>`

View pull request details.

**Output fields:** ID, Title, State, From Branch, To Branch, Author, Reviewers (with approval status), Created, Updated, URL, Description

### `bitbucket pr create <project-key> <repo-slug>`

Create a pull request.

| Flag | Required | Description |
|------|----------|-------------|
| `--title <text>` | Yes | PR title |
| `--from <branch>` | Yes | Source branch |
| `--to <branch>` | Yes | Target branch |
| `--description <text>` | No | PR description |
| `--reviewers <slugs>` | No | Comma-separated reviewer usernames |

```
orbit -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add new feature" \
  --reviewers john.doe,jane.smith
```

### `bitbucket pr edit <project-key> <repo-slug> <pr-id>`

Edit a pull request's title or description. **Alias:** `update`. **Server only.**

| Flag | Required | Description |
|------|----------|-------------|
| `--title <text>` | No | New title |
| `--description <text>` | No | New description |
| `--description-file <path>` | No | Read the description from a file; `-` reads stdin. Mutually exclusive with `--description` |

At least one of the three is required. Use `--description-file` for anything longer than a line:

```
orbit -p myprofile bb pr edit L3SUP agents-sre 42 --title "Correct the title"
orbit -p myprofile bb pr edit L3SUP agents-sre 42 --description-file pr-body.md
cat pr-body.md | orbit -p myprofile bb pr edit L3SUP agents-sre 42 --description-file -
```

Reviewers survive the edit. Bitbucket's update endpoint replaces the pull request instead of patching it, so a request sending only a description also means "and no reviewers" - it drops every reviewer and their approvals and still answers 200. This command sends the existing reviewers and unchanged title back with the edit, then verifies they survived; a dropped reviewer or a reset approval is reported as an error naming them, since only that person can approve again. Do not hand-roll this with `pr` writes elsewhere.

### `bitbucket pr diff <project-key> <repo-slug> <pr-id>`

Show the unified diff for a pull request. Useful for code review.

| Flag | Default | Description |
|------|---------|-------------|
| `--context <n>` | 3 | Number of context lines around changes |

```
orbit -p myprofile bb pr diff L3SUP agents-sre 42
orbit -p myprofile bb pr diff L3SUP agents-sre 42 --context 10
```

### `bitbucket pr merge <project-key> <repo-slug> <pr-id>`

Merge a pull request. Automatically handles version for optimistic locking.

| Flag | Default | Description |
|------|---------|-------------|
| `--bypass-review` | `false` | Temporarily disable PRE_PULL_REQUEST_MERGE hooks AND set project-level default reviewer conditions to 0 required approvals before merging, then restore all settings after. Requires project/repo admin permissions. |

```
orbit -p myprofile bb pr merge L3SUP agents-sre 42
orbit -p myprofile bb pr merge L3SUP agents-sre 42 --bypass-review
```

### `bitbucket pr merge-status <project-key> <repo-slug> <pr-id>`

Whether a pull request can be merged, and what is blocking it. **Aliases:** `can-merge`, `mergeable`. **Server only.**

| Flag | Default | Description |
|------|---------|-------------|
| `--exit-status` | `false` | Exit non-zero when it cannot be merged |

```
orbit -p myprofile bb pr merge-status L3SUP agents-sre 42
orbit -p myprofile bb pr merge-status L3SUP agents-sre 42 --exit-status
```

Shows the server's verdict (`CLEAN`, `CONFLICTED`, `VETOED`), the conflict or each vetoing merge check in its own words, the reviewer tally with positions first, and the head commit's build stats. Ask this before `pr merge` rather than reading a merge failure.

An already merged or declined PR is reported as such: Bitbucket answers the merge question with a 409 for one that is not open, so the state is read first and becomes the answer. Build statuses come from a separate API that an instance need not have; when it is absent the line says so and the verdict still stands.

### `bitbucket pr approve <project-key> <repo-slug> <pr-id>`

Approve a pull request as the current user.

### `bitbucket pr unapprove <project-key> <repo-slug> <pr-id>`

Remove your approval from a pull request.

### `bitbucket pr needs-work <project-key> <repo-slug> <pr-id>`

**Alias:** `request-changes`

Flag a pull request as `NEEDS_WORK` - Bitbucket Server's "request changes" signal. Blocks the merge without closing the PR (unlike `decline`). **Server only.** Pair it with `pr comment` so the author has the actionable feedback.

All three status commands verify the status Bitbucket stored, so a success line means the flag is really on the PR.

### `bitbucket pr decline <project-key> <repo-slug> <pr-id>`

Decline a pull request. Automatically handles version for optimistic locking.

### `bitbucket pr comment <project-key> <repo-slug> <pr-id>`

Add a comment to a pull request.

| Flag | Required | Description |
|------|----------|-------------|
| `--body <text>` | Yes | Comment text |

### `bitbucket pr activity <project-key> <repo-slug> <pr-id>`

List pull request activity — comments, approvals, merges, status changes.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

**Activity types shown:** COMMENTED, APPROVED, REVIEWED (needs work), MERGED, DECLINED, OPENED

Nested comment replies are shown with indentation — all reply threads are traversed recursively so no responses are hidden.

---

## code

### `bitbucket code search <query>`

Search file contents and paths across repositories. **Aliases:** `find`, `grep`. **Server only**, needs the search index.

| Flag | Default | Description |
|------|---------|-------------|
| `--project <key>` | | One project |
| `--repo <slug>` | | One repository (requires `--project`) |
| `--repos` | `false` | List matching repositories instead of every file |
| `--limit <n>` | 25 | Max matching files |

**This is how you tell what a repo is built with** - names rarely say:

```
orbit -p myprofile bb code search go.mod --repos             # Go
orbit -p myprofile bb code search build.gradle.kts --repos   # Kotlin/Gradle
orbit -p myprofile bb code search pom.xml --repos            # Maven
orbit -p myprofile bb code search "TODO" --project EPCAP --repo agentic-platform-monorepo
```

**Only `project:` and `repo:` modifiers exist, and an unsupported one is silent.** `lang:kotlin` returns zero hits, exactly like a search that matched nothing - do not conclude from it that there is no Kotlin. Use a build file instead. `--project`/`--repo` set the real modifiers; any other `word:` in the query is flagged in the output.

---

## user

### `bitbucket user list`

List users.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter by username or display name |
| `--limit <n>` | 25 | Max results |

---

## reviewer-condition (admin)

Alias: `rc`

Manage project-level default reviewer conditions. These conditions auto-assign reviewers to PRs and can enforce a minimum number of required approvals before merging.

### `bitbucket reviewer-condition list <project-key>`

List all default reviewer conditions for a project. Shows condition ID, scope, source/target branch matchers, required approvals count, and reviewer names.

```
orbit -p myprofile bb rc list EPCAP
orbit -p myprofile bb rc list EPCAP -o json
```

**Output fields:** ID, Scope (PROJECT/REPOSITORY), Source, Target, Required Approvals, Reviewers

### `bitbucket reviewer-condition update <project-key> <condition-id>`

Update the required approvals count for a default reviewer condition. Preserves all reviewers and branch matchers — only changes the approval requirement.

| Flag | Required | Description |
|------|----------|-------------|
| `--required-approvals <n>` | Yes | Number of required approvals (0 = no requirement) |

```
# Temporarily bypass required approvals
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 0

# Restore to 2 required approvals
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 2
```

**Common use case — bypass merge block:**
When `pr merge` fails with "Not all required reviewers have approved yet", this is caused by a project-level condition. Set required approvals to 0, merge, then restore. The `--bypass-review` flag on `pr merge` does this automatically.

### `bitbucket reviewer-condition delete <project-key> <condition-id>`

Delete a default reviewer condition permanently.

```
orbit -p myprofile bb rc delete EPCAP 1063
```
