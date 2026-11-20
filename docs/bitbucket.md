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
- [Code Search](#code)
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
| `--limit` | int | - | Maximum number of projects to return. |

```bash
# List all projects
orbit bb project list -p myprofile

# List with a limit
orbit bb project list -p myprofile --limit 10
```

### `project search`

Search projects whose **display name** matches the text. **Alias:** `find`. **Server/Data Center only.**

| Argument | Position | Description |
|----------|----------|-------------|
| `text` | 1 | Part of a project's name. Optional when `--permission` is given. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--permission` | string | | Keep only projects you hold this on: `PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`. |
| `--limit` | int | 25 | Max results. |

```bash
orbit bb project search agentic -p myprofile
orbit bb project search "Ad Server" -p myprofile
orbit bb project search "" --permission PROJECT_ADMIN --limit 100 -p myprofile
```

> The match is on the **name**, not the key. A project keyed `EPCAP` and named "Epsilon PeopleCloud Agentic Platform" does not turn up for `EPCAP`; use `project view EPCAP` when you already have the key.

### `project view`

View details of a single project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key (e.g. `TEAM`). |

```bash
orbit bb project view TEAM -p myprofile
```

### `project permissions`

List who holds a project: each user's and each group's own grant (`PROJECT_READ`, `PROJECT_WRITE`, `PROJECT_ADMIN`, `REPO_CREATE`), and what the project gives **every licensed user** (`default`). Bitbucket shows this only to the project's admins; anyone else gets "only an admin of project KEY can read who holds it".

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

```bash
orbit bb project permissions EPCAP
orbit bb project permissions EPCAP -o json   # {users, groups, default}
```

> The default is easy to miss and usually the answer. On EPCAP every licensed user holds `PROJECT_READ`, which is why someone with no grant of their own anywhere still reads every repository in it.

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
| `--limit` | int | - | Maximum number of repositories to return. |

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

### `repo search`

Search repositories by name **across every project** you can see. **Alias:** `find`. **Server/Data Center only.**

| Argument | Position | Description |
|----------|----------|-------------|
| `text` | 1 | Part of a repository name. Optional when `--project` or `--project-name` is given. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | | Restrict to one project, **by key**. |
| `--project-name` | string | | Restrict to projects whose display name matches this text. Mutually exclusive with `--project`. |
| `--permission` | string | | Keep only repos you hold this on: `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`. |
| `--visibility` | string | | `public` or `private`. |
| `--archived` | string | | `ACTIVE` (the server's default), `ARCHIVED` or `ALL`. |
| `--limit` | int | 25 | Max results. |

```bash
# Across every project - no project key needed, unlike `repo list`
orbit bb repo search agentic -p myprofile

# Narrowed to one project, by key
orbit bb repo search agentic --project EPCAP -p myprofile

# Only the ones you could push to
orbit bb repo search payments --permission REPO_WRITE -p myprofile

# Every repo in a family of projects
orbit bb repo search "" --project-name "Ad Server" --limit 100 -p myprofile
```

> **`--project` takes a key; Bitbucket's own filter does not.** The API filters repositories by `projectname`, which matches a project's *display name*: asking it for `projectname=EPCAP` returns an empty page, with a `200` and no error, on an instance whose EPCAP project is named "Epsilon PeopleCloud Agentic Platform". Since a key is what every other command here takes, `--project` resolves the key to that project's name to narrow the search server-side, then matches the results on the key itself - a name fragment can belong to several projects ("Acce" matches both "Access" and "Access Business Components").
>
> Paging follows the server's own `nextPageStart` rather than `start + limit`. Bitbucket filters each page by permission *after* reading it, so the offset it hands back routinely jumps past what it returned - a three-result page on a live instance came back with `nextPageStart=34`. Computing the offset would re-read some repositories and silently skip others.

### `repo create`

Create a new repository in a project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | *(required)* | Repository name. |
| `--description` | string | | Repository description. |
| `--default-branch` | string | | Default branch name (e.g. `main`). |
| `--public` | bool | `false` | Make the repository publicly visible. |
| `--forkable` | bool | `true` | Allow forking. |

```bash
orbit bb repo create L3SUP --name my-service --description "Payment gateway" -p myprofile
orbit bb repo create L3SUP --name my-api --public -p myprofile
```

### `repo edit`

Update a repository's settings. **Alias:** `repo update`.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | New repository name. |
| `--description` | string | New description. |
| `--default-branch` | string | New default branch. |
| `--public` | bool | Set public visibility. |
| `--forkable` | bool | Allow/disallow forking. |

```bash
orbit bb repo edit L3SUP my-service --description "Updated description" -p myprofile
orbit bb repo edit L3SUP my-service --public=false --forkable=false -p myprofile
```

### `repo delete`

Delete a repository. **Irreversible.**

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

```bash
orbit bb repo delete L3SUP my-service -p myprofile
```

### `repo fork`

Fork a repository, optionally into a different project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | Source project key. |
| `repo-slug` | 2 | Source repository slug. |

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | Fork name (default: same as source). |
| `--target-project` | string | Target project key (default: same project). |

```bash
orbit bb repo fork L3SUP agents-sre --name agents-sre-fork -p myprofile
orbit bb repo fork L3SUP agents-sre --target-project MYPROJ -p myprofile
```

### `repo forks`

List forks of a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 25 | Maximum forks to return. |

```bash
orbit bb repo forks L3SUP agents-sre -p myprofile
```

### `repo merge-strategy`

Show or set which merge strategies a repository allows and its default. **Server only.** Setting needs repo admin, makes the repository use its own configuration (rather than its project's or the server's), and reads it back.

```
orbit bb repo merge-strategy <project-key> <repo-slug> [--default <id> --enable <id,id>]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--default` | string | | Strategy a merge uses unless it asks for another |
| `--enable` | strings | | Strategies a merge may use; must include the default |

Strategy ids: `no-ff`, `ff`, `ff-only`, `rebase-no-ff`, `rebase-ff-only`, `squash`, `squash-ff-only`.

```bash
orbit bb repo merge-strategy EPCAP agentic-platform-monorepo -p myprofile
orbit bb repo merge-strategy EPCAP agentic-platform-monorepo --default squash --enable squash,no-ff -p myprofile
```

### `repo permissions`

List the users holding their own grant on a repository; `--groups` adds the groups holding one. Neither says what one person can do on it - they may hold it through the project, a group or the project's default. That is [`repo access`](#repo-access).

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Maximum results. |
| `--groups` | bool | false | Also list group grants. The JSON is then `{users, groups}`; without it, the user list as before. |

```bash
orbit bb repo permissions L3SUP agents-sre -p myprofile
orbit bb repo permissions L3SUP agents-sre --groups -o json
```

### `repo access`

What one person can do on a repository, and every grant it comes from: their own grant on the repository, a group's grant on it, their own or a group's grant on its project, the project's grant to every licensed user, and a public repository. The person is one account exactly, by username or email (see [`user view`](#user-view)).

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `username-or-email` | 3 | The account's username or email address, exactly. |

```bash
orbit bb repo access EPCAP agentic-platform-chat-sdk karthik.reddy1@epsilon.com
# Karthik Reddy (karreddy8@publicisgroupe.net) on EPCAP/agentic-platform-chat-sdk: REPO_WRITE
#   REPO_WRITE   from repository (REPO_WRITE)
#   REPO_READ    from project default (PROJECT_READ)
```

JSON (`-o json`): `user`, `effective` (the highest `REPO_*`, `""` for none), `direct` (their own grant on the repository - the one `repo grant` changes), `sources` (`via`, `permission` as Bitbucket holds it, `level` on this repository), `membership_known` and `unresolved_groups`.

> Group membership is read through the endpoint Bitbucket's own user pickers use. Where an instance hides it from the account orbit signs in as, `membership_known` is false, the groups holding a grant are listed as unresolved, and `effective` is a lower bound ("at least ..."). Reading the grants needs admin on the project and the repository, as `project permissions` does.

### `repo grant`

Grant a user permission on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `username` | 3 | Bitbucket username (slug). |

| Flag | Type | Description |
|------|------|-------------|
| `--permission` | string | *(required)* One of `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`. |

```bash
orbit bb repo grant L3SUP agents-sre john.doe --permission REPO_WRITE -p myprofile
```

### `repo revoke`

Revoke a user's permission on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `username` | 3 | Bitbucket username (slug). |

```bash
orbit bb repo revoke L3SUP agents-sre john.doe -p myprofile
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
| `--limit` | int | - | Maximum number of branches to return. |

```bash
orbit bb branch list TEAM my-service -p myprofile
orbit bb branch list TEAM my-service -p myprofile --limit 20
```

### `branch create`

Create a new branch.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `name` | 3 | Name of the new branch. |
| `start-point` | 4 | The branch name, tag or commit SHA to branch from. |

```bash
# Create a feature branch from main
orbit bb branch create TEAM my-service feature/login main -p myprofile

# Create a branch from a specific commit
orbit bb branch create TEAM my-service hotfix/urgent a1b2c3d -p myprofile
```

### `branch delete`

Delete a branch. **Server/Data Center only** - Bitbucket Cloud keeps branches under a different API.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `branch` | 3 | Name of the branch to delete. Either form works: `feature/login` or `refs/heads/feature/login`. |

```bash
orbit bb branch delete TEAM my-service feature/login -p myprofile
```

The commit the branch was at is printed, because that is the way back - recreate it with `branch create` from that SHA:

```
Deleted branch scratch/orbit-final-check (was at 5e85a5908712)
```

> **Two things about this operation are not what they appear, and both are handled here.**
>
> Branch deletion is not in Bitbucket's core REST API. `/rest/api/latest/.../branches` serves `GET` only, so a `DELETE` sent there comes back as `405 Method Not Allowed` with an empty body - no branch name, no reason, the same answer for every branch. Deletion lives in the `branch-utils` plugin, a separate REST namespace rather than a different version of the same one. Orbit did send the `DELETE` to the core API until v0.72.0, which is why every delete failed with a 405.
>
> And `branch-utils` answers a delete of a branch that **is not there** with the same `204 No Content` it gives one it did delete. A typo would therefore report success and leave the branch in place, which is worse than the 405: someone clearing out a list of branches would have no way to tell which names were real. So the branch is looked up first and looked for again afterwards, and "deleted" is only printed when it was there and then was not:
>
> ```
> $ orbit bb branch delete TEAM my-service feature/typo
> Error: no branch feature/typo in TEAM/my-service
> ```
>
> Bitbucket's own refusals come through intact, since they say what to do about them - for example `You cannot delete this branch since it has at least one open pull request associated with it`.

### `branch default`

Show the default branch of a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

```bash
orbit bb branch default TEAM my-service -p myprofile
```

To **change** the default branch, use [`repo edit`](#repo-edit):

```bash
orbit bb repo edit TEAM my-service --default-branch develop -p myprofile
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
| `--limit` | int | - | Maximum number of tags to return. |

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
| `tag-name` | 3 | Name of the tag. |
| `start-point` | 4 | The branch name or commit SHA to tag. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `-m`, `--message` | string | No | Tag message. Creates an annotated tag. |

```bash
# Tag the current HEAD of main
orbit bb tag create TEAM my-service v1.0.0 main -p myprofile

# Tag a specific commit, annotated
orbit bb tag create TEAM my-service v1.0.1 a1b2c3d -m "Release v1.0.1" -p myprofile
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
| `--limit` | int | - | Maximum number of commits to return. |

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
| `--state` | string | - | Filter by state: `open`, `declined`, `merged`, or `all`. |
| `--limit` | int | 25 | Maximum number of pull requests to return. |
| `--draft` | bool | - | Keep only drafts. `--draft=false` keeps only pull requests that are not drafts. Bitbucket filters server side. |

```bash
# List open PRs (default)
orbit bb pr list TEAM my-service -p myprofile

# List merged PRs
orbit bb pr list TEAM my-service -p myprofile --state merged

# List all PRs with a limit
orbit bb pr ls TEAM my-service -p myprofile --state all --limit 25

# Only the drafts, and only the ones that are not
orbit bb pr list TEAM my-service --draft -p myprofile
orbit bb pr list TEAM my-service --draft=false -p myprofile
```

The `STATE` column reads `DRAFT` for an open draft, the way the web interface shows it: a draft is `OPEN`, but a listing that called it that would hide the only unusual thing about it. A closed pull request keeps whatever draft flag it had - Bitbucket does not clear it - so there the column shows `MERGED` or `DECLINED`, which is the fact that matters.

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

The `State` line folds in Bitbucket's separate draft flag: `OPEN (draft)` is a pull request nobody is being asked to merge yet. See [`pr draft`](#pr-draft--pr-ready).

### `pr create`

Create a new pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--from` | string | Yes | Source branch name. |
| `--to` | string | Yes | Target branch name. |
| `--title` | string | Yes | Pull request title. |
| `--description` | string | No | Pull request description body. |
| `--reviewers` | string | No | Comma-separated reviewer usernames. |
| `--draft` | bool | No | Create as a draft. Bitbucket will not merge it until it is marked ready. Requires Data Center 8.18 or newer. |

```bash
# Minimal PR
orbit bb pr create TEAM my-service \
  --from feature/login \
  --to main \
  --title "Add login flow" \
  -p myprofile

# PR with description and reviewers
orbit bb pr create TEAM my-service \
  --from feature/login \
  --to main \
  --title "Add login flow" \
  --description "Implements OAuth2 login with refresh token support." \
  --reviewers john.doe,jane.smith \
  -p myprofile

# A draft, for work that is not ready for review yet
orbit bb pr create TEAM my-service \
  --from feature/login \
  --to main \
  --title "Add login flow" \
  --draft \
  -p myprofile
```

> **A draft that was not made is reported.** An instance older than 8.18 does not know the field, ignores it, and answers `201` with an ordinary pull request - one that is open for review when it was meant to be held back. orbit checks what came back and says so, naming the pull request it did create.

### `pr edit`

Edit the title, description, target branch or draft state of an open pull request. **Alias:** `update`. **Server/Data Center only** - Bitbucket Cloud has a different pull request API.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Description |
|------|------|-------------|
| `--title` | string | New title. |
| `--description` | string | New description. |
| `--description-file` | string | Read the new description from a file; `-` reads stdin. Mutually exclusive with `--description`. |
| `--to` | string | New target branch - the branch the pull request will merge into. For several pull requests at once, see [`pr retarget`](#pr-retarget). |
| `--draft` | bool | Convert to a draft; `--draft=false` marks it ready for review. For several at once, see [`pr draft`](#pr-draft--pr-ready). |

At least one of them is required.

A description is usually longer than fits comfortably on a command line, which is what `--description-file` is for:

```bash
# A one-line title fix
orbit bb pr edit TEAM my-service 42 --title "Add login flow (OAuth2)" -p myprofile

# A real description, from a file
orbit bb pr edit TEAM my-service 42 --description-file pr-body.md -p myprofile

# ...or from a pipe
generate-release-notes | orbit bb pr edit TEAM my-service 42 --description-file - -p myprofile

# A new target branch, on its own or along with the text, in one update
orbit bb pr edit TEAM my-service 42 --to develop -p myprofile

# A retitle that also pulls it back to a draft, in one update
orbit bb pr edit TEAM my-service 42 --title "WIP: spike" --draft -p myprofile
```

> **Reviewers are preserved, and that takes work.** Bitbucket's update endpoint *replaces* the pull request rather than patching it: a request carrying only a new description also says "and no reviewers", so Bitbucket drops every reviewer and every approval they had given - and answers `200` as if nothing happened. Editing a description would silently reset the review.
>
> So this command reads the pull request first (it needs the `version` for optimistic locking anyway), sends the existing reviewers and the unchanged title back along with the edit, and then checks what came back. A reviewer who disappeared or an approval that was reset is reported as an error naming them, because neither is something orbit can restore - only the reviewer can approve again.
>
> Two differences between what was sent and what Bitbucket stores are the server normalising its input, not a failed edit, and are accepted: CRLF line endings become LF, and trailing whitespace is trimmed (a description read from a file ends in a newline; the stored copy does not).

### `pr retarget`

Point one or more open pull requests at a different target branch - what the web interface calls editing the destination branch. **Server/Data Center only.**

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id...` | 3+ | One or more pull request IDs. A leading `#` is accepted. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--to` | string | Yes | The branch the pull requests will merge into. A branch name or a full `refs/heads/...` ref. |

```bash
# One pull request
orbit bb pr retarget TEAM my-service 42 --to develop -p myprofile

# A stack of them, in one go
orbit bb pr retarget TEAM my-service 42 43 44 --to develop -p myprofile

# Machine-readable, including which ones were already there
orbit bb pr retarget TEAM my-service 42 43 44 --to develop -o json -p myprofile

# Every open pull request that targets main, moved onto develop
orbit bb pr list TEAM my-service --state open -o json -p myprofile \
  | jq -r '.[] | select(.toRef.displayId == "main") | .id' \
  | xargs orbit bb pr retarget TEAM my-service --to develop -p myprofile
```

```
PR #42 (feature/one): main -> develop
PR #43 (feature/two): already targets develop
PR #44: it is MERGED, not open; Bitbucket would move the target branch of a closed pull request anyway, and it would still never merge
1 retargeted, 1 already there, 1 failed
```

A run over several pull requests does not stop at the first failure: each is reported on its own line, the tally follows, and the command exits non-zero if any of them did not move. Running it again is safe - the ones that moved are left alone the second time.

> **What this command refuses to do, and why.**
>
> The target branch is checked once, before anything is written. A mistyped branch name would otherwise be found by the first pull request and reported again for every pull request after it, with the earlier ones already moved.
>
> **Reviewers and approvals are preserved.** The update endpoint replaces the pull request rather than patching it, so an update carrying only the new branch also says "and no reviewers" - see the note under [`pr edit`](#pr-edit). A retarget reads the pull request first, sends the reviewers back with the new branch, and reports any approval the update cost.
>
> **A pull request that already targets the branch is left alone.** Bitbucket accepts that update happily, but it is not free: the version is bumped and everyone watching is notified. Retargeting a stack one branch at a time would otherwise send mail about work that did not happen.
>
> **A merged or declined pull request is reported, not moved.** Bitbucket will change the target branch of a closed pull request without complaint, leaving a pull request that points at a branch it never merged into. Over a list of pull requests that is a silent rewrite of the record.
>
> **Only one pull request may be open for a given source and target.** Bitbucket answers a retarget that would collide with a `409` whose body is the entire conflicting pull request - some four kilobytes of JSON around one sentence that does not include its number. orbit reports the number instead: `#7 is already open from feature/x to develop, and Bitbucket allows only one pull request for a given source and target; merge, decline or retarget #7 first`.

### `pr draft` / `pr ready`

Convert one or more open pull requests to drafts, or mark drafts ready for review - the web interface's "Convert to draft" and "Mark as ready". **Requires Bitbucket Data Center 8.18 or newer. Server/Data Center only.**

`pr ready` has the aliases `publish` and `undraft`.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id...` | 3+ | One or more pull request IDs. A leading `#` is accepted. |

Neither takes a flag; the direction is the command.

```bash
# Hold one back
orbit bb pr draft TEAM my-service 42 -p myprofile

# Send a whole stack out for review
orbit bb pr ready TEAM my-service 42 43 44 -p myprofile

# Machine-readable, including which ones were already there
orbit bb pr draft TEAM my-service 42 43 44 -o json -p myprofile

# Every draft in the repository, marked ready
orbit bb pr list TEAM my-service --draft -o json -p myprofile \
  | jq -r '.[].id' \
  | xargs orbit bb pr ready TEAM my-service -p myprofile
```

```
PR #42 (feature/one): now a draft
PR #43 (feature/two): already a draft
PR #44: it is MERGED, not open, and Bitbucket will not change the draft state of a closed pull request - it would have to be reopened first
1 marked as draft, 1 already in draft, 1 failed
```

As with `pr retarget`, a run over several pull requests does not stop at the first failure, and running it again is safe.

> **What a draft is, in Bitbucket's terms.** It is a flag on the pull request, not a state of its own: a draft is `OPEN`, it can be reviewed and approved, and approvals given to a draft survive it being marked ready. The one thing it cannot do is merge - the merge refuses with *"This is a draft pull request and cannot be merged"*, and `orbit bb pr merge` refuses it before that, pointing at `pr ready`.
>
> **The merge check does not know about it.** `orbit bb pr merge-status` asks the server whether a pull request can merge, and for a draft the server answers `CLEAN`, with no conflict and no veto. orbit counts the draft as a blocker itself, so `merge-status` reports `Mergeable: no` with the draft as the reason rather than promising a merge that would be refused.
>
> **Reviewers and approvals are preserved.** The draft flag rides on the same endpoint that replaces the pull request, so an update carrying only the flag would drop every reviewer - see the note under [`pr edit`](#pr-edit). The pull request is read first, the reviewers go back with the change, and any approval it cost is reported.
>
> **A pull request already in the state asked for is left alone.** Bitbucket accepts that update, bumps the version and notifies everyone watching, for a change that did not happen.
>
> **A merged or declined pull request is reported, not changed.** Bitbucket refuses this one itself - *"This pull request has been declined and must be reopened before it can be drafted or published"* - but only after the write has been built and sent. orbit refuses it first, which is what names the pull request in the message when several were asked for at once.
>
> **An instance without drafts is caught.** Older instances ignore the field and answer `200` with the pull request unchanged. orbit checks the flag that came back rather than the status code, and reports the version that has drafts.

### `pr merge`

Merge a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--bypass-review` | bool | `false` | Temporarily disable merge checks and required-approval conditions for the merge, restoring them afterwards. Requires repo admin. |
| `--strategy` | string | repo default | Merge strategy id (`squash`, `no-ff`, ...). Must be one the repository enables; checked before anything else, including `--bypass-review`. With `--strategy` the new commit's parent count is checked against it (1 for a squash, 2 for a merge commit); a merge that happened but cannot be confirmed is still reported as merged, follow-up steps still run, and the command exits non-zero. Without it nothing is checked. |
| `--message` | string | | The whole commit message, subject line first. |
| `--message-from-pr` | bool | `false` | Commit message `Pull request #N: title` with the description as the body - what the web UI writes. |
| `--delete-source-branch` | bool | `false` | Delete the source branch once merged, only while it still points at the merged commit. Left in place, and reported, when it is in a fork (or the response does not say), is the default branch or a branch-model development/production branch, has another open pull request into it (the base of a stack), or was pushed to since. |

```bash
orbit bb pr merge TEAM my-service 42 -p myprofile
orbit bb pr merge TEAM my-service 42 --strategy squash --message-from-pr --delete-source-branch -p myprofile
```

> **The commit message over the API.** A merge made through the REST API, with no message, gets Bitbucket's subject `Merge pull request #N in P/R from X to Y` - the title is not in it, unlike a merge from the web UI (`Pull request #N: title`). `--message-from-pr` writes the UI's form, so history reads the same however a pull request was merged. For a squash, Bitbucket then appends `Squashed commit of the following:` and every branch commit, whatever the message; neither the API nor the repository settings turn that off.

> **A draft is refused before anything is touched.** Bitbucket will not merge a draft, and it says so only once the merge has been attempted - which under `--bypass-review` means every merge hook disabled and every required-approvals condition zeroed first, all for a request that was never going to be accepted. So the pull request's draft flag is checked up front:
>
> ```
> pull request TEAM/my-service#42 is a draft, and Bitbucket will not merge a draft; mark it ready for review first:
>   orbit bb pr ready TEAM my-service 42
> ```
>
> **A refusal reads as one sentence.** Bitbucket answers a refused merge with a `409` whose body is a JSON object wrapped around the reason; the reason is kept and the wrapper dropped, so an already-merged pull request reports *"This pull request has already been merged."* rather than the exception that carried it. A pull request that changed between being read and being merged says which version each side was on, and that nothing was merged.

### `pr merge-status`

Show whether a pull request can be merged, and what is blocking it. **Aliases:** `can-merge`, `mergeable`. **Server/Data Center only.**

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--exit-status` | bool | `false` | Exit non-zero when the pull request cannot be merged. |

```bash
orbit bb pr merge-status TEAM my-service 42 -p myprofile

# Gate a script on it without parsing output
orbit bb pr merge-status TEAM my-service 42 --exit-status -p myprofile
```

It reports the server's own verdict - the one the merge button acts on - and the reasons behind it:

```
PR:        #1298  fix(EAP-1832): tighten catalog OpenAPI types
State:     OPEN  (fix/EAP-1832-catalog-any-waiver -> develop)
Mergeable: no (CONFLICTED)
Blocked by:
  - the source branch conflicts with the target and has to be rebased or merged first
Reviewers: 1 approved, 0 needs work, 0 no response (of 1)
  Sailesh Rajagopalan          APPROVED
Builds:    none reported
URL:       https://git.example.com/projects/TEAM/repos/my-service/pull-requests/1298
```

`Blocked by` carries the conflict, if any, and then each merge check that vetoed the merge with its own wording - "Not enough approvals", "Requires successful builds", or whatever an add-on hook says. Reviewers who have taken a position are listed first, so one approval among a dozen silent reviewers is visible at a glance.

> **A draft is counted as a blocker, because the server does not count it.** Bitbucket's merge check answers `CLEAN` for a draft - no conflict, no veto - and then the merge endpoint refuses it. Relaying only the server's verdict would promise a merge that cannot happen, so the draft is added as a blocker of its own:
>
> ```
> State:     OPEN (draft)  (feature/login -> main)
> Mergeable: no
> Blocked by:
>   - the pull request is a draft, and Bitbucket will not merge a draft until it is marked ready for review
> ```
>
> The outcome in brackets is the server's merge check, not always what decided the verdict, so it is left out when it would contradict it: `no (CLEAN)` would read as a contradiction, and the blocker already says what is wrong. It is printed as usual when the two agree - the `no (CONFLICTED)` above, or a `yes (CLEAN)` for one that can merge. `--exit-status` and the `canMerge` field of `-o json` follow the verdict, so a script gated on either will not try to merge a draft.

> **An already merged or declined pull request is reported as such, not as an error.** Bitbucket answers the merge question with a `409` for one that is not open ("This pull request has already been merged."), so the pull request's state is read first and becomes the answer:
>
> ```
> State:     MERGED  (fix/EAP-2360-agent-slug-v2 -> develop)
> Mergeable: no - the pull request is merged, so there is no merge to make
> ```
>
> Build statuses come from `/rest/build-status`, a separate API that an instance need not have. When it cannot be read the line says so and the verdict still stands - a missing build-status plugin is not a reason to fail a question about approvals.

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

Add a comment to a pull request, or with `--file` and `--line` to one line of its diff, as a reviewer leaves it in the web UI (Server and Data Center).

`--line` is the line's number in the new file; `--old-line` its number in the old file, for a line the pull request removes. Whether the line is added, removed or unchanged context is read from the pull request's own diff, with the web UI's 10 lines of context, so the caller never states it. A line the diff does not show is refused (`not in the pull request's diff`): it has nowhere to anchor.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--body` | string | Yes | The comment text. |
| `--file` | string | No | Comment on a line of this file's diff (its path in the repository). |
| `--line` | int | No | With `--file`: the line's number in the new file. |
| `--old-line` | int | No | With `--file`: the line's number in the old file, for a removed line. |
| `--blocker` | bool | No | With `--file`: make it a task the author has to resolve. |
| `--commit` | string | No | With `--file`: refuse if the pull request's source has moved past this commit (a prefix is enough), so a comment written against one head never lands on another. |

```bash
orbit bb pr comment TEAM my-service 42 --body "LGTM, approved." -p myprofile
orbit bb pr comment EPCAP agentic-platform-chat-sdk 6 --file gocd/scripts/publish.sh --line 12 \
  --body "This runs before the version check" --commit 0395020f
orbit bb pr comment EPCAP agentic-platform-chat-sdk 6 --file README.md --old-line 630 --blocker \
  --body "Keep the VPN note: the publish needs it" -o json
```

### `pr dashboard`

Your pull requests across every repository, as your Bitbucket dashboard lists them (Server and Data Center). `--role reviewer` is what waits on your review.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--role` | string | No | `reviewer`, `author` or `participant` (default: any). |
| `--state` | string | No | `OPEN` (default), `MERGED` or `DECLINED`. |
| `--limit` | int | No | Max results (default: every one). |

```bash
orbit bb pr dashboard --role reviewer
orbit bb pr dashboard --role author --state merged --limit 10 -o json
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

### `pr needs-work`

Flag a pull request as `NEEDS_WORK` - Bitbucket Server's "request changes" signal. Blocks merge the same way the manual "Needs Work" button in the UI does, without closing the PR (unlike `decline`). Aliases: `request-changes`. **Server only** - Bitbucket Cloud has a different approval model.

Works even when the caller isn't already listed as a reviewer on the PR, and fails loudly if Bitbucket doesn't actually store the flag - see the note under this section.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
# Leave the review comment first, then flip to needs-work
orbit bb pr comment TEAM my-service 42 -m "Please address the three items above" -p myprofile
orbit bb pr needs-work TEAM my-service 42 -p myprofile

# When the author pushes fixes, approve to clear the status
orbit bb pr approve TEAM my-service 42 -p myprofile
```

> `approve`, `unapprove`, and `needs-work` all share a single `SetReviewStatus` helper that `PUT`s to `/participants/{userSlug}` with the target status, then checks the participant Bitbucket echoes back to confirm the status was stored. The check is not paranoia: `POST /participants` accepts a `status` field, answers `200`, and stores only the role - so the earlier POST-based implementation reported success while leaving the caller `UNAPPROVED` and the PR unflagged.
>
> If the `PUT` is refused (some instances reject it for a user who isn't on the PR yet, though Bitbucket documents it as adding them implicitly), the helper `POST`s to `/participants` with `role=REVIEWER` to enrol the caller and retries the `PUT` once. If that fails too, the error you see is Bitbucket's own refusal of the status change.
>
> The caller's username is resolved from the profile's `auth.username` (basic auth) or from the `X-AUSERNAME` header on an authed GET (PAT/bearer), then memoized for the life of the process.

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

Nested comment replies are shown with indentation - all reply threads are traversed recursively.

---

## code

Search code across repositories. **Server/Data Center only**, and it needs the instance's search index.

### `code search`

Search file contents and paths. **Aliases:** `find`, `grep`.

| Argument | Position | Description |
|----------|----------|-------------|
| `query` | 1 | Bitbucket search syntax. A bare term matches file contents *and* paths. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | | Restrict to one project, by key. |
| `--repo` | string | | Restrict to one repository. Requires `--project`. |
| `--repos` | bool | `false` | List the repositories that match instead of every file. |
| `--limit` | int | 25 | Max matching files to fetch. |

**This is how you find what a repository is built with**, which its name rarely says:

```bash
orbit bb code search go.mod --repos -p myprofile              # Go modules
orbit bb code search build.gradle.kts --repos -p myprofile    # Kotlin / Gradle
orbit bb code search pom.xml --repos -p myprofile             # Maven
orbit bb code search requirements.txt --repos -p myprofile    # Python

# Scoped
orbit bb code search go.mod --project EPCAP -p myprofile
orbit bb code search "TODO" --project EPCAP --repo agentic-platform-monorepo -p myprofile
```

`--repos` collapses the file matches into the list of repositories holding them, with a count each, which is usually the actual question:

```
PROJECT          REPOSITORY                     FILES
-------          ----------                     -----
DEVOPSES         consul-k8s                     3
EPCAP            agentic-platform-monorepo      1
HDP              gometrics                      1
```

> **Only `project:` and `repo:` are supported modifiers, and an unsupported one is not an error.** `lang:kotlin` returns zero hits - identical to a search that genuinely matched nothing - so a query written with `lang:` or `ext:` reads as "this instance has no Kotlin". `--project` and `--repo` set the two real modifiers for you, and any other `word:` in a query is called out in the output for exactly this reason. Bitbucket also rejects `repo:` without a `project:`, which this command catches before spending a request.

---

## user

Manage Bitbucket users.

### `user list`

List users.

**Alias:** `user ls`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--filter` | string | - | Filter by username or display name. |
| `--limit` | int | 25 | Maximum number of users to return. |

```bash
orbit bb user list -p myprofile
orbit bb user ls -p myprofile --limit 100
orbit bb user list --filter john -p myprofile
```

**Output columns:** ID, USERNAME, DISPLAY NAME, ACTIVE. The columns are sized from the data rather than to a fixed width, because on an instance where the username is an email address it runs well past twenty characters - and it cannot be truncated, since it is the value that `pr create --reviewers` takes.

> `--filter` is a substring search over username, display name **and** email. "Karthik Reddy" and "Karthik Meyyappan" both match `karthik`; to name one person, use `user view`.

### `user view`

The one account whose username or email address is exactly the one given (case-insensitive), across every page of the directory. Refused when no account has it, and when two accounts share the email (give the username then). A display name is not accepted: two people can share one.

| Argument | Position | Description |
|----------|----------|-------------|
| `username-or-email` | 1 | The username or email address, exactly. |

```bash
orbit bb user view karthik.reddy1@epsilon.com
orbit bb user view karreddy8@publicisgroupe.net -o json
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
