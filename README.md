# orbit

A unified CLI for managing connections to development lifecycle services.

Supports **Jira**, **Confluence**, **GitLab**, **GitHub**, and **Bitbucket** (cloud and self-hosted). Organize connections into profiles to switch between projects seamlessly.

Secrets can be stored as [1Password](https://1password.com/) references (`op://vault/item/field`) and are resolved at runtime using the 1Password CLI.

## Install

### macOS — Homebrew

```bash
brew install jorgemuza/tap/orbit
```

### Windows — Scoop

```bash
scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket
scoop install orbit
```

### From source

```bash
go install github.com/jorgemuza/orbit@latest
```

### Binary releases

Download from [GitHub Releases](https://github.com/jorgemuza/orbit/releases).

## Quick Start

```bash
# Create a profile
orbit profile create --name my-project --default

# Add services
orbit profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method basic --username me@co.com --token "op://Dev/jira-token/credential"

orbit profile add-service \
  --name gitlab-onprem --type gitlab --variant server \
  --base-url https://gitlab.internal.com \
  --auth-method basic --username admin --password "op://Dev/gitlab/password"

orbit profile add-service \
  --name github-cloud --type github \
  --auth-method token --token "op://Dev/github-token/credential"

# Test connectivity
orbit service ping
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `~/.config/orbit/config.yaml` |
| `-p, --profile` | Profile to use (overrides default) | |
| `-o, --output` | Output format: `table`, `json`, `yaml` | `table` |

## Commands Overview

### Profile & Service Management

```bash
orbit profile create --name proj --default         # Create a profile
orbit profile list                                  # List profiles
orbit profile show proj                             # Show profile details
orbit profile use proj                              # Set default profile
orbit profile delete proj                           # Delete profile
orbit profile add-service --name svc --type jira ...  # Add service
orbit profile remove-service svc                    # Remove service

orbit service list                                  # List services in active profile
orbit service ping                                  # Test all connections
orbit service ping jira-cloud                       # Test specific service

orbit auth                                          # Pre-resolve 1Password secrets
orbit auth clear                                    # Clear cached secrets

orbit version                                       # Print version info
```

### Jira

Manage issues, epics, sprints, boards, projects, releases, custom fields, screens, statuses, and issue types. Supports Jira Cloud (API v3) and Server/Data Center (API v2).

```bash
orbit jira issue list --project PROJ --assignee me
orbit jira issue create --project PROJ --type Story --summary "Add login"
orbit jira issue edit PROJ-123 --field customfield_10397="Yes"
orbit jira issue move PROJ-123 "In Progress"
orbit jira epic list --project PROJ
orbit jira sprint list --board-id 1 --state active
orbit jira field list --custom
orbit jira field create --name "AI Assisted" --type select
orbit jira screen list --filter "PYMT"
orbit jira screen tab-create 10089 "AI Workflow"
orbit jira screen field-add 10089 10868 --fields "customfield_10397,customfield_10398"
orbit jira status list
orbit jira issuetype-list
```

**[Full Jira reference →](docs/jira.md)**

### Confluence

Manage pages, publish markdown directories, and control page layout. Supports Cloud and Server/Data Center.

```bash
orbit confluence page 12345
orbit confluence children 12345
orbit confluence create --space DEV --parent 12345 --title "Guide" --file docs/guide.md
orbit confluence update 12345 --file docs/guide.md
orbit confluence publish ./docs --space DEV --parent 12345 --dry-run
orbit confluence set-width 12345 --recursive
```

**[Full Confluence reference →](docs/confluence.md)**

### GitLab

Manage projects, merge requests, pipelines, jobs, branches, tags, issues, runners, variables, and more. Alias: `gl`.

```bash
orbit gl project view my-team/app
orbit gl project create --name "my-app" --namespace my-team
orbit gl projects --group my-team
orbit gl mr list my-team/app --state opened
orbit gl mr create my-team/app --source feature/x --target main --title "Add feature"
orbit gl pipeline list my-team/app
orbit gl pipeline jobs my-team/app 12345
orbit gl job log my-team/app 67890
orbit gl branch list my-team/app
orbit gl variable list my-team/app
```

**[Full GitLab reference →](docs/gitlab.md)**

### GitHub

Manage repositories, pull requests, issues, Actions workflow runs, releases, secrets, and more. Alias: `gh`.

```bash
orbit gh repos --org my-org
orbit gh pr list octocat/hello-world
orbit gh pr create octocat/hello-world --head feature/x --base main --title "Add feature"
orbit gh run list octocat/hello-world
orbit gh run watch octocat/hello-world 12345
orbit gh issue list octocat/hello-world --labels bug
orbit gh secret set octocat/hello-world MY_SECRET --value "s3cret"
orbit gh release latest octocat/hello-world
```

**[Full GitHub reference →](docs/github.md)**

### Bitbucket

Manage projects, repositories, pull requests, branches, tags, and pipelines. Supports Cloud and Data Center. Alias: `bb`.

```bash
orbit bb repo list MY-PROJ
orbit bb pr list MY-PROJ my-repo --state open
orbit bb pr create MY-PROJ my-repo --source feature/x --target main --title "Add feature"
orbit bb branch list MY-PROJ my-repo
```

**[Full Bitbucket reference →](docs/bitbucket.md)**

## Supported Services

| Service | Type | Variants | Default Base URL |
|---------|------|----------|-----------------|
| Jira | `jira` | `cloud`, `server` | *(required)* |
| Confluence | `confluence` | `cloud`, `server` | *(required)* |
| GitLab | `gitlab` | `cloud`, `server` | `https://gitlab.com` |
| GitHub | `github` | `cloud`, `server` | `https://api.github.com` |
| Bitbucket | `bitbucket` | `cloud`, `server` | `https://api.bitbucket.org/2.0` |

### Authentication Methods

| Method | Flags | Use Case |
|--------|-------|----------|
| `token` | `--token` | API tokens, PATs (most common) |
| `basic` | `--username`, `--password` | Basic auth (Jira Cloud uses email + API token) |
| `oauth2` | `--client-id`, `--client-secret` | OAuth2 client credentials |

## 1Password Integration

Instead of storing secrets in plain text, use 1Password references:

```bash
orbit profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method basic \
  --username me@company.com \
  --password "op://DevVault/jira-token/credential"
```

Secrets are resolved at runtime via `op read`, so the [1Password CLI](https://developer.1password.com/docs/cli/) must be installed and authenticated. Resolved values are cached locally to avoid repeated biometric prompts.

Pre-resolve secrets before running multiple commands:

```bash
orbit auth          # triggers biometric once, caches secrets
orbit auth clear    # remove cached secrets
```

By default, cached secrets expire after **8 hours**. You can change this (or disable expiration entirely) in `config.yaml`:

```yaml
settings:
  secrets_cache_ttl_hours: 0   # 0 = never expire (default: 8)
```

## Configuration

Config is stored in YAML at `~/.config/orbit/config.yaml`:

```yaml
settings:
  secrets_cache_ttl_hours: 8   # 0 = never expire (default: 8)

profiles:
  - name: my-project
    description: "My project services"
    default: true
    services:
      - name: jira-cloud
        type: jira
        variant: cloud
        base_url: https://myco.atlassian.net
        auth:
          method: basic
          username: me@company.com
          password: "op://DevVault/jira-token/credential"
      - name: gitlab-onprem
        type: gitlab
        variant: server
        base_url: https://gitlab.internal.com
        auth:
          method: basic
          username: admin
          password: "op://DevVault/gitlab/password"
      - name: github-cloud
        type: github
        auth:
          method: token
          token: "op://DevVault/github-token/credential"
```

## Claude Code Skills

Orbit ships with [Claude Code](https://docs.anthropic.com/en/docs/claude-code) skills that give Claude deep knowledge of each service's CLI commands, flags, and workflows. Install them with `npx @anthropic-ai/claude-code-skills`:

```bash
# Install all orbit skills at once
npx @anthropic-ai/claude-code-skills --skills jira,confluence,github,gitlab,bitbucket,format-docs --from github:jorgemuza/orbit

# Install a single skill
npx @anthropic-ai/claude-code-skills --skills jira --from github:jorgemuza/orbit

# Install a subset
npx @anthropic-ai/claude-code-skills --skills jira,confluence --from github:jorgemuza/orbit
```

Skills are installed into `.claude/skills/` in your current project directory. Once installed, Claude Code automatically discovers them and can use orbit commands on your behalf.

### Available Skills

| Skill | Description |
|-------|-------------|
| `jira` | Issue CRUD, epics, sprints, boards, fields, screens, statuses, wiki markup formatting |
| `confluence` | Page view/create/update, markdown publishing, page width control |
| `github` | Repos, PRs, Actions runs, issues, releases, secrets (alias: `gh`) |
| `gitlab` | Projects, MRs, pipelines, branches, tags, variables (alias: `gl`) |
| `bitbucket` | Projects, repos, PRs, branches, tags, pipelines (alias: `bb`) |
| `format-docs` | Prepare and restructure markdown for Confluence publishing |

## Releasing

Releases are automated via GitHub Actions. Push a tag to trigger `goreleaser`:

```bash
git tag -a v0.4.0 -m "v0.4.0: description"
git push origin v0.4.0
```

> **Warning:** Do not run `goreleaser release` locally and then push the tag. The CI will fail because the release assets already exist. Always let CI handle the release by only pushing the tag.

## License

MIT
