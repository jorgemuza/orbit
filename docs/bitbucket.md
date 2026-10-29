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
| `--limit` | int | — | Maximum number of projects to return. |

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

### `repo permissions`

List user-level permissions on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Maximum results. |

```bash
orbit bb repo permissions L3SUP agents-sre -p myprofile
```

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
| `--limit` | int | — | Maximum number of branches to return. |

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
| `--from` | string | Yes | Source branch name. |
| `--to` | string | Yes | Target branch name. |
| `--title` | string | Yes | Pull request title. |
| `--description` | string | No | Pull request description body. |
| `--reviewers` | string | No | Comma-separated reviewer usernames. |

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
```

### `pr edit`

Edit the title or description of an open pull request. **Alias:** `update`. **Server/Data Center only** - Bitbucket Cloud has a different pull request API.

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

A description is usually longer than fits comfortably on a command line, which is what `--description-file` is for:

```bash
# A one-line title fix
orbit bb pr edit TEAM my-service 42 --title "Add login flow (OAuth2)" -p myprofile

# A real description, from a file
orbit bb pr edit TEAM my-service 42 --description-file pr-body.md -p myprofile

# ...or from a pipe
generate-release-notes | orbit bb pr edit TEAM my-service 42 --description-file - -p myprofile
```

> **Reviewers are preserved, and that takes work.** Bitbucket's update endpoint *replaces* the pull request rather than patching it: a request carrying only a new description also says "and no reviewers", so Bitbucket drops every reviewer and every approval they had given - and answers `200` as if nothing happened. Editing a description would silently reset the review.
>
> So this command reads the pull request first (it needs the `version` for optimistic locking anyway), sends the existing reviewers and the unchanged title back along with the edit, and then checks what came back. A reviewer who disappeared or an approval that was reset is reported as an error naming them, because neither is something orbit can restore - only the reviewer can approve again.
>
> Two differences between what was sent and what Bitbucket stores are the server normalising its input, not a failed edit, and are accepted: CRLF line endings become LF, and trailing whitespace is trimmed (a description read from a file ends in a newline; the stored copy does not).

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

### `pr needs-work`

Flag a pull request as `NEEDS_WORK` — Bitbucket Server's "request changes" signal. Blocks merge the same way the manual "Needs Work" button in the UI does, without closing the PR (unlike `decline`). Aliases: `request-changes`. **Server only** — Bitbucket Cloud has a different approval model.

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

Nested comment replies are shown with indentation — all reply threads are traversed recursively.

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
