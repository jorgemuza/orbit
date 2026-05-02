# Orbit GoCD Command Reference

Manage GoCD pipelines, agents, environments, config repos, server administration, security, elastic agents, and more from the command line.

**Top-level command:** `orbit gocd` (alias: `cd`)

**Persistent flag (all subcommands):**

| Flag | Description |
|------|-------------|
| `--service` | GoCD service name, required only when a profile has multiple GoCD services configured |

**Notes:**

- GoCD is always self-hosted -- `base_url` is required in the service configuration.
- Authentication uses either token (Bearer) or basic auth.
- All commands support `-o json` and `-o yaml` for structured output.
- **GoCD 25.x compatible.** All API versions verified against [api.gocd.org/25.1.0](https://api.gocd.org/25.1.0/). Dashboard handles both v4 (25.x) and older response formats. Pipeline instance falls back to history extraction when the direct instance API is restricted by permission.

---

## Table of Contents

- [pipeline -- Pipeline commands](#pipeline)
  - [pipeline list](#pipeline-list)
  - [pipeline status](#pipeline-status)
  - [pipeline history](#pipeline-history)
  - [pipeline get](#pipeline-get)
  - [pipeline trigger](#pipeline-trigger)
  - [pipeline pause](#pipeline-pause)
  - [pipeline unpause](#pipeline-unpause)
  - [pipeline create](#pipeline-create)
  - [pipeline update](#pipeline-update)
  - [pipeline delete](#pipeline-delete)
  - [pipeline comment](#pipeline-comment)
  - [pipeline export](#pipeline-export)
  - [pipeline compare](#pipeline-compare)
  - [pipeline lock](#pipeline-lock)
  - [pipeline unlock](#pipeline-unlock)
- [pipeline-group -- Pipeline group commands](#pipeline-group)
  - [pipeline-group list](#pipeline-group-list)
  - [pipeline-group get](#pipeline-group-get)
  - [pipeline-group create](#pipeline-group-create)
  - [pipeline-group update](#pipeline-group-update)
  - [pipeline-group delete](#pipeline-group-delete)
- [access-token -- Access token commands](#access-token)
  - [access-token list](#access-token-list)
  - [access-token create](#access-token-create)
  - [access-token revoke](#access-token-revoke)
  - [access-token list-all](#access-token-list-all)
  - [access-token revoke-admin](#access-token-revoke-admin)
- [agent -- Agent commands](#agent)
  - [agent list](#agent-list)
  - [agent get](#agent-get)
  - [agent enable](#agent-enable)
  - [agent disable](#agent-disable)
  - [agent delete](#agent-delete)
  - [agent kill-task](#agent-kill-task)
  - [agent update](#agent-update)
  - [agent job-history](#agent-job-history)
- [environment -- Environment commands](#environment)
  - [environment list](#environment-list)
  - [environment get](#environment-get)
  - [environment create](#environment-create)
  - [environment delete](#environment-delete)
  - [environment update](#environment-update)
  - [environment patch](#environment-patch)
- [config-repo -- Config repo commands](#config-repo)
  - [config-repo list](#config-repo-list)
  - [config-repo get](#config-repo-get)
  - [config-repo status](#config-repo-status)
  - [config-repo trigger](#config-repo-trigger)
  - [config-repo create](#config-repo-create)
  - [config-repo update](#config-repo-update)
  - [config-repo delete](#config-repo-delete)
  - [config-repo definitions](#config-repo-definitions)
  - [config-repo preflight](#config-repo-preflight)
- [user -- User commands](#user)
  - [user list](#user-list)
  - [user get](#user-get)
  - [user create](#user-create)
  - [user update](#user-update)
  - [user delete](#user-delete)
  - [user delete-bulk](#user-delete-bulk)
  - [user current](#user-current)
  - [user update-current](#user-update-current)
- [plugin -- Plugin commands](#plugin)
  - [plugin list](#plugin-list)
  - [plugin get](#plugin-get)
  - [plugin get-settings](#plugin-get-settings)
  - [plugin create-settings](#plugin-create-settings)
  - [plugin update-settings](#plugin-update-settings)
- [role -- Role commands](#role)
  - [role list](#role-list)
  - [role get](#role-get)
  - [role create](#role-create)
  - [role update](#role-update)
  - [role delete](#role-delete)
- [authorization -- Authorization commands](#authorization)
  - [authorization list](#authorization-list)
  - [authorization get](#authorization-get)
  - [authorization create](#authorization-create)
  - [authorization update](#authorization-update)
  - [authorization delete](#authorization-delete)
- [backup -- Backup commands](#backup)
  - [backup get-config](#backup-get-config)
  - [backup create-config](#backup-create-config)
  - [backup delete-config](#backup-delete-config)
  - [backup schedule](#backup-schedule)
- [cluster-profile -- Cluster profile commands](#cluster-profile)
  - [cluster-profile list](#cluster-profile-list)
  - [cluster-profile get](#cluster-profile-get)
  - [cluster-profile create](#cluster-profile-create)
  - [cluster-profile update](#cluster-profile-update)
  - [cluster-profile delete](#cluster-profile-delete)
- [dashboard -- Dashboard](#dashboard)
- [elastic-agent-profile -- Elastic agent profile commands](#elastic-agent-profile)
  - [elastic-agent-profile list](#elastic-agent-profile-list)
  - [elastic-agent-profile get](#elastic-agent-profile-get)
  - [elastic-agent-profile create](#elastic-agent-profile-create)
  - [elastic-agent-profile update](#elastic-agent-profile-update)
  - [elastic-agent-profile delete](#elastic-agent-profile-delete)
  - [elastic-agent-profile usage](#elastic-agent-profile-usage)
- [material -- Material commands](#material)
  - [material list](#material-list)
  - [material get](#material-get)
  - [material usage](#material-usage)
  - [material trigger-update](#material-trigger-update)
- [notification-filter -- Notification filter commands](#notification-filter)
  - [notification-filter list](#notification-filter-list)
  - [notification-filter get](#notification-filter-get)
  - [notification-filter create](#notification-filter-create)
  - [notification-filter delete](#notification-filter-delete)
- [package -- Package commands](#package)
  - [package list](#package-list)
  - [package get](#package-get)
  - [package create](#package-create)
  - [package update](#package-update)
  - [package delete](#package-delete)
  - [package usage](#package-usage)
- [package-repo -- Package repository commands](#package-repo)
  - [package-repo list](#package-repo-list)
  - [package-repo get](#package-repo-get)
  - [package-repo create](#package-repo-create)
  - [package-repo update](#package-repo-update)
  - [package-repo delete](#package-repo-delete)
- [artifact -- Artifact store commands](#artifact)
  - [artifact list-store](#artifact-list-store)
  - [artifact get-store](#artifact-get-store)
  - [artifact create-store](#artifact-create-store)
  - [artifact update-store](#artifact-update-store)
  - [artifact delete-store](#artifact-delete-store)
- [secret-config -- Secret configuration commands](#secret-config)
  - [secret-config list](#secret-config-list)
  - [secret-config get](#secret-config-get)
  - [secret-config create](#secret-config-create)
  - [secret-config update](#secret-config-update)
  - [secret-config delete](#secret-config-delete)
- [stage -- Stage commands](#stage)
  - [stage cancel](#stage-cancel)
  - [stage run](#stage-run)
- [job -- Job commands](#job)
  - [job run](#job-run)
- [server -- Server commands](#server)
  - [server health](#server-health)
  - [server maintenance](#server-maintenance)
  - [server maintenance-on](#server-maintenance-on)
  - [server maintenance-off](#server-maintenance-off)
- [server-config -- Server configuration commands](#server-config)
  - [server-config site-url get](#server-config-site-url-get)
  - [server-config site-url update](#server-config-site-url-update)
  - [server-config artifact-config get](#server-config-artifact-config-get)
  - [server-config artifact-config update](#server-config-artifact-config-update)
  - [server-config job-timeout get](#server-config-job-timeout-get)
  - [server-config job-timeout update](#server-config-job-timeout-update)
  - [server-config mail-server get](#server-config-mail-server-get)
  - [server-config mail-server update](#server-config-mail-server-update)
  - [server-config mail-server delete](#server-config-mail-server-delete)
- [template -- Template commands](#template)
  - [template list](#template-list)
  - [template get](#template-get)
  - [template create](#template-create)
  - [template update](#template-update)
  - [template delete](#template-delete)
- [version -- Server version](#version)
- [encrypt -- Encrypt a value](#encrypt)

---

## pipeline

Manage GoCD pipelines.

### pipeline list

List all pipelines grouped by pipeline group.

```
orbit gocd pipeline list [flags]
```

**Example:**

```bash
orbit gocd pipeline list
orbit cd pipeline list -o json
```

### pipeline status

Show pipeline status (paused, schedulable, locked).

```
orbit gocd pipeline status <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

**Example:**

```bash
orbit gocd pipeline status my-pipeline
orbit cd pipeline status my-pipeline -o json
```

### pipeline history

Show pipeline run history.

```
orbit gocd pipeline history <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `10` | Maximum number of history entries |

**Example:**

```bash
orbit gocd pipeline history my-pipeline
orbit cd pipeline history my-pipeline --limit 5 -o json
```

### pipeline get

Get a pipeline instance (latest or by counter). Uses a three-strategy fallback for maximum GoCD version compatibility: modern instance path (GoCD 18.x+), legacy instance path (GoCD 14.3–17.x), then extraction from pipeline history (works on all versions, including servers that restrict the instance endpoint by permission — confirmed on GoCD 25.1.0).

```
orbit gocd pipeline get <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--counter` | `0` | Pipeline counter (0 = latest from history) |

**Example:**

```bash
orbit gocd pipeline get my-pipeline
orbit cd pipeline get my-pipeline --counter 42 -o json
```

### pipeline trigger

Schedule a pipeline run with optional per-run environment variable and material revision overrides.

```
orbit gocd pipeline trigger <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Type | Description |
|------|------|-------------|
| `--env` | string (repeatable) | Environment variable override (`KEY=VALUE`). Overrides the pipeline's default for this run only. |
| `--secure-env` | string (repeatable) | Mark an env var name as secure (pair with `--env`). The value is encrypted server-side. |
| `--material` | string (repeatable) | Pin a material to a specific revision (`FINGERPRINT=REVISION`). |

Without any flags, trigger behaves as before (empty body, uses pipeline defaults).

**Examples:**

```bash
# Simple trigger (no overrides)
orbit gocd pipeline trigger my-pipeline

# Override env vars for this run
orbit cd pipeline trigger my-pipeline --env VERSION=1.2.3 --env REGION=us-east-1

# Trigger with a secure env var (value encrypted server-side)
orbit cd pipeline trigger my-pipeline --env "TF_VAR_db_password=s3cret" --secure-env TF_VAR_db_password

# Pin a material to a specific revision
orbit cd pipeline trigger my-pipeline --material abc123=a2d23c5

# Secret rotation pattern: trigger with per-run overrides pointing at
# the right encrypted var and AWS Secrets Manager secret ID
orbit cd pipeline trigger d1epap-secrets-v2 \
  --env SECRET_ID=epsilon/dev/phoenix-otel-headers \
  --env SECRET_KEY=SECRET_PHOENIX_OTEL_HEADERS
```

### pipeline pause

Pause a pipeline.

```
orbit gocd pipeline pause <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--reason` | `""` | Reason for pausing |

**Example:**

```bash
orbit gocd pipeline pause my-pipeline --reason "maintenance"
orbit cd pipeline pause my-pipeline
```

### pipeline unpause

Unpause a pipeline.

```
orbit gocd pipeline unpause <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

**Example:**

```bash
orbit gocd pipeline unpause my-pipeline
```

### pipeline create

Create a pipeline from a JSON or YAML file.

```
orbit gocd pipeline create --group <group> --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--group` | Yes | Pipeline group name |
| `--from-file` | Yes | Path to JSON or YAML file with pipeline definition |

**Example:**

```bash
orbit gocd pipeline create --group my-group --from-file pipeline.yaml
orbit cd pipeline create --group my-group --from-file pipeline.json
```

### pipeline update

Update a pipeline configuration from a file. Automatically fetches the current ETag.

```
orbit gocd pipeline update <name> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline configuration |

> **Secure env vars preserved.** Round-trip `config -o json > file.json` → `update --from-file file.json` is safe for pipelines with secure environment variables (e.g. `TF_VAR_db_password`). The `encrypted_value` field is preserved on unmarshal and sent back with `omitempty`, so GoCD keeps the existing encrypted blob. Only the populated field (`value` for plaintext, `encrypted_value` for secure) is included in the update payload.

**Example:**

```bash
# Safe round-trip: export then re-import (secure vars preserved)
orbit cd pipeline config my-pipeline -o json > pipeline.json
# edit pipeline.json (add a stage, change a material, etc.)
orbit cd pipeline update my-pipeline --from-file pipeline.json

orbit gocd pipeline update my-pipeline --from-file pipeline.yaml
```

### pipeline delete

Delete a pipeline.

```
orbit gocd pipeline delete <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

**Example:**

```bash
orbit gocd pipeline delete my-pipeline
```

### pipeline comment

Add a comment to a pipeline instance.

```
orbit gocd pipeline comment <name> --counter <N> --message <msg>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Required | Description |
|------|----------|-------------|
| `--counter` | Yes | Pipeline counter |
| `--message` | Yes | Comment message |

**Example:**

```bash
orbit gocd pipeline comment my-pipeline --counter 42 --message "Approved for prod"
orbit cd pipeline comment my-pipeline --counter 42 --message "Hotfix deployed"
```

### pipeline export

Export pipeline configuration. Optionally specify a config repo plugin.

```
orbit gocd pipeline export <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--plugin-id` | `""` | Config repo plugin ID (optional) |

**Example:**

```bash
orbit gocd pipeline export my-pipeline
orbit cd pipeline export my-pipeline --plugin-id yaml.config.plugin -o json
```

### pipeline compare

Compare two pipeline instances to see material changes between them.

```
orbit gocd pipeline compare <name> --from <N> --to <N> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from` | Yes | From pipeline counter |
| `--to` | Yes | To pipeline counter |

**Example:**

```bash
orbit gocd pipeline compare my-pipeline --from 1 --to 2
orbit cd pipeline compare my-pipeline --from 5 --to 10 -o json
```

### pipeline lock

Lock a pipeline to prevent further scheduling until it is unlocked.

```
orbit gocd pipeline lock <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

**Example:**

```bash
orbit gocd pipeline lock my-pipeline
```

### pipeline unlock

Unlock a locked pipeline.

```
orbit gocd pipeline unlock <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

**Example:**

```bash
orbit gocd pipeline unlock my-pipeline
```

---

## pipeline-group

Manage GoCD pipeline groups.

**Aliases:** `pg`

### pipeline-group list

List all pipeline groups.

```
orbit gocd pipeline-group list [flags]
```

**Example:**

```bash
orbit gocd pipeline-group list
orbit cd pg list -o json
```

### pipeline-group get

Get pipeline group details.

```
orbit gocd pipeline-group get <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline group name |

**Example:**

```bash
orbit gocd pipeline-group get my-group
orbit cd pg get my-group -o json
```

### pipeline-group create

Create a pipeline group from a file.

```
orbit gocd pipeline-group create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline group definition |

**Example:**

```bash
orbit gocd pipeline-group create --from-file group.yaml
orbit cd pg create --from-file group.json
```

### pipeline-group update

Update a pipeline group. Automatically fetches the current ETag.

```
orbit gocd pipeline-group update <name> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline group name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline group configuration |

**Example:**

```bash
orbit gocd pipeline-group update my-group --from-file group.yaml
orbit cd pg update my-group --from-file group.json
```

### pipeline-group delete

Delete a pipeline group.

```
orbit gocd pipeline-group delete <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Pipeline group name |

**Example:**

```bash
orbit gocd pipeline-group delete my-group
orbit cd pg delete my-group
```

---

## agent

Manage GoCD agents.

### agent list

List all agents.

```
orbit gocd agent list [flags]
```

**Example:**

```bash
orbit gocd agent list
orbit cd agent list -o json
```

### agent get

View agent details.

```
orbit gocd agent get <uuid>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent get adb9540a-b954-4571-9d9b-2f330739d4da
orbit cd agent get abc-123-uuid -o json
```

### agent enable

Enable an agent. Uses ETag-based optimistic locking.

```
orbit gocd agent enable <uuid>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent enable adb9540a-b954-4571-9d9b-2f330739d4da
```

### agent disable

Disable an agent. Uses ETag-based optimistic locking.

```
orbit gocd agent disable <uuid>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent disable adb9540a-b954-4571-9d9b-2f330739d4da
```

### agent delete

Delete an agent.

```
orbit gocd agent delete <uuid>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent delete adb9540a-b954-4571-9d9b-2f330739d4da
```

### agent kill-task

Kill running tasks on an agent.

```
orbit gocd agent kill-task <uuid>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent kill-task adb9540a-b954-4571-9d9b-2f330739d4da
```

### agent update

Update an agent from a file.

```
orbit gocd agent update <uuid> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with agent update data |

**Example:**

```bash
orbit gocd agent update abc-123-uuid --from-file agent.yaml
orbit cd agent update abc-123-uuid --from-file agent.json
```

### agent job-history

Show job run history for an agent.

```
orbit gocd agent job-history <uuid> [flags]
```

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

**Example:**

```bash
orbit gocd agent job-history abc-123-uuid
orbit cd agent job-history abc-123-uuid -o json
```

---

## environment

Manage GoCD environments.

**Aliases:** `env`

### environment list

List all environments.

```
orbit gocd environment list [flags]
```

**Example:**

```bash
orbit gocd environment list
orbit cd env list -o json
```

### environment get

View environment details including pipelines, agents, and environment variables.

```
orbit gocd environment get <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

**Example:**

```bash
orbit gocd environment get production
orbit cd env get staging -o json
```

### environment create

Create a new environment.

```
orbit gocd environment create <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

| Flag | Description |
|------|-------------|
| `--pipeline` | Pipeline to add (repeatable) |

**Example:**

```bash
orbit gocd environment create staging
orbit cd env create staging --pipeline my-app --pipeline my-api
```

### environment delete

Delete an environment.

```
orbit gocd environment delete <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

**Example:**

```bash
orbit gocd environment delete staging
```

### environment update

Replace an environment definition from a file. Automatically fetches the current ETag.

```
orbit gocd environment update <name> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with environment definition |

**Example:**

```bash
orbit gocd environment update staging --from-file env.yaml
orbit cd env update staging --from-file env.json
```

### environment patch

Patch an environment by adding or removing pipelines and agents.

```
orbit gocd environment patch <name> [flags]
```

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

| Flag | Description |
|------|-------------|
| `--add-pipeline` | Pipeline to add (repeatable) |
| `--remove-pipeline` | Pipeline to remove (repeatable) |
| `--add-agent` | Agent to add (repeatable) |
| `--remove-agent` | Agent to remove (repeatable) |

**Example:**

```bash
orbit gocd environment patch staging --add-pipeline p1 --remove-pipeline p2
orbit cd env patch staging --add-agent a1 --remove-agent a2
```

---

## config-repo

Manage GoCD config repositories.

**Aliases:** `configrepo`, `cr`

### config-repo list

List config repos.

```
orbit gocd config-repo list [flags]
```

**Example:**

```bash
orbit gocd config-repo list
orbit cd cr list -o json
```

### config-repo get

View config repo details.

```
orbit gocd config-repo get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

**Example:**

```bash
orbit gocd config-repo get my-config-repo
orbit cd cr get my-repo -o json
```

### config-repo status

Check config repo sync status.

```
orbit gocd config-repo status <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

**Example:**

```bash
orbit gocd config-repo status my-config-repo
orbit cd cr status my-repo -o json
```

### config-repo trigger

Trigger a config repo update.

```
orbit gocd config-repo trigger <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

**Example:**

```bash
orbit gocd config-repo trigger my-config-repo
orbit cd cr trigger my-repo
```

### config-repo create

Create a config repository from a file.

```
orbit gocd config-repo create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with config repo definition |

**Example:**

```bash
orbit gocd config-repo create --from-file repo.yaml
orbit cd cr create --from-file repo.json
```

### config-repo update

Update a config repository from a file. Automatically fetches the current ETag.

```
orbit gocd config-repo update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with config repo definition |

**Example:**

```bash
orbit gocd config-repo update my-repo --from-file repo.yaml
orbit cd cr update my-repo --from-file repo.json
```

### config-repo delete

Delete a config repository.

```
orbit gocd config-repo delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

**Example:**

```bash
orbit gocd config-repo delete my-repo
orbit cd cr delete my-repo
```

### config-repo definitions

Get the pipelines and environments defined by a config repo.

```
orbit gocd config-repo definitions <id> [flags]
```

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

**Example:**

```bash
orbit gocd config-repo definitions my-repo
orbit cd cr definitions my-repo -o json
```

### config-repo preflight

Run a preflight check on config repo content before applying it.

```
orbit gocd config-repo preflight --plugin-id <id> --from-file <file> [flags]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--plugin-id` | Yes | Config repo plugin ID |
| `--from-file` | Yes | Path to JSON or YAML file with content to check |
| `--repo-id` | No | Config repo ID (optional, for existing repos) |

**Example:**

```bash
orbit gocd config-repo preflight --plugin-id yaml.config.plugin --from-file pipeline.gocd.yaml
orbit cd cr preflight --plugin-id yaml.config.plugin --repo-id my-repo --from-file pipeline.gocd.yaml
```

---

## user

Manage GoCD users.

### user list

List all users.

```
orbit gocd user list [flags]
```

**Example:**

```bash
orbit gocd user list
orbit cd user list -o json
```

### user get

Get user details.

```
orbit gocd user get <login>
```

| Argument | Description |
|----------|-------------|
| `login` | User login name |

**Example:**

```bash
orbit gocd user get jdoe
orbit cd user get jdoe -o json
```

### user create

Create a user from a file.

```
orbit gocd user create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with user definition |

**Example:**

```bash
orbit gocd user create --from-file user.yaml
orbit cd user create --from-file user.json
```

### user update

Update a user from a file.

```
orbit gocd user update <login> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `login` | User login name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with user configuration |

**Example:**

```bash
orbit gocd user update jdoe --from-file user.yaml
orbit cd user update jdoe --from-file user.json
```

### user delete

Delete a user.

```
orbit gocd user delete <login>
```

| Argument | Description |
|----------|-------------|
| `login` | User login name |

**Example:**

```bash
orbit gocd user delete jdoe
```

### user delete-bulk

Delete multiple users at once.

```
orbit gocd user delete-bulk [flags]
```

| Flag | Description |
|------|-------------|
| `--user` | User login name to delete (repeatable) |
| `--from-file` | Path to JSON or YAML file with bulk delete request |

**Example:**

```bash
orbit gocd user delete-bulk --user jdoe --user jane
orbit cd user delete-bulk --from-file users.yaml
```

### user current

Show the current authenticated user.

```
orbit gocd user current [flags]
```

**Example:**

```bash
orbit gocd user current
orbit cd user current -o json
```

### user update-current

Update the current authenticated user from a file.

```
orbit gocd user update-current --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with user fields to update |

**Example:**

```bash
orbit gocd user update-current --from-file user.yaml
orbit cd user update-current --from-file user.json
```

---

## plugin

Manage GoCD plugins.

### plugin list

List installed plugins.

```
orbit gocd plugin list [flags]
```

**Example:**

```bash
orbit gocd plugin list
orbit cd plugin list -o json
```

### plugin get

Get plugin details.

```
orbit gocd plugin get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Plugin ID |

**Example:**

```bash
orbit gocd plugin get my-plugin
orbit cd plugin get my-plugin -o json
```

### plugin get-settings

Get plugin settings.

```
orbit gocd plugin get-settings <plugin-id>
```

| Argument | Description |
|----------|-------------|
| `plugin-id` | Plugin ID |

**Example:**

```bash
orbit gocd plugin get-settings my-plugin
orbit cd plugin get-settings my-plugin -o json
```

### plugin create-settings

Create plugin settings from a file.

```
orbit gocd plugin create-settings --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with plugin settings |

**Example:**

```bash
orbit gocd plugin create-settings --from-file settings.yaml
orbit cd plugin create-settings --from-file settings.json
```

### plugin update-settings

Update plugin settings from a file. Automatically fetches the current ETag.

```
orbit gocd plugin update-settings <plugin-id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `plugin-id` | Plugin ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with plugin settings |

**Example:**

```bash
orbit gocd plugin update-settings my-plugin --from-file settings.yaml
orbit cd plugin update-settings my-plugin --from-file settings.json
```

---

## role

Manage GoCD security roles.

### role list

List all roles.

```
orbit gocd role list [flags]
```

**Example:**

```bash
orbit gocd role list
orbit cd role list -o json
```

### role get

Get role details.

```
orbit gocd role get <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Role name |

**Example:**

```bash
orbit gocd role get my-role
orbit cd role get my-role -o json
```

### role create

Create a role from a file.

```
orbit gocd role create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with role definition |

**Example:**

```bash
orbit gocd role create --from-file role.yaml
orbit cd role create --from-file role.json
```

### role update

Update a role from a file. Automatically fetches the current ETag.

```
orbit gocd role update <name> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `name` | Role name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with role configuration |

**Example:**

```bash
orbit gocd role update my-role --from-file role.yaml
orbit cd role update my-role --from-file role.json
```

### role delete

Delete a role.

```
orbit gocd role delete <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Role name |

**Example:**

```bash
orbit gocd role delete my-role
```

---

## authorization

Manage GoCD authorization configurations.

**Aliases:** `auth-config`

### authorization list

List authorization configurations.

```
orbit gocd authorization list [flags]
```

**Example:**

```bash
orbit gocd authorization list
orbit cd auth-config list -o json
```

### authorization get

Get authorization configuration details.

```
orbit gocd authorization get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Authorization configuration ID |

**Example:**

```bash
orbit gocd authorization get my-auth
orbit cd auth-config get my-auth -o json
```

### authorization create

Create an authorization configuration from a file.

```
orbit gocd authorization create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with authorization configuration |

**Example:**

```bash
orbit gocd authorization create --from-file auth.yaml
orbit cd auth-config create --from-file auth.json
```

### authorization update

Update an authorization configuration from a file. Automatically fetches the current ETag.

```
orbit gocd authorization update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Authorization configuration ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with authorization configuration |

**Example:**

```bash
orbit gocd authorization update my-auth --from-file auth.yaml
orbit cd auth-config update my-auth --from-file auth.json
```

### authorization delete

Delete an authorization configuration.

```
orbit gocd authorization delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Authorization configuration ID |

**Example:**

```bash
orbit gocd authorization delete my-auth
orbit cd auth-config delete my-auth
```

---

## backup

Manage GoCD server backups.

### backup get-config

Get backup configuration.

```
orbit gocd backup get-config [flags]
```

**Example:**

```bash
orbit gocd backup get-config
orbit cd backup get-config -o json
```

### backup create-config

Create backup configuration from a file.

```
orbit gocd backup create-config --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with backup configuration |

**Example:**

```bash
orbit gocd backup create-config --from-file backup.yaml
orbit cd backup create-config --from-file backup.json
```

### backup delete-config

Delete backup configuration.

```
orbit gocd backup delete-config
```

**Example:**

```bash
orbit gocd backup delete-config
```

### backup schedule

Schedule a server backup.

```
orbit gocd backup schedule [flags]
```

**Example:**

```bash
orbit gocd backup schedule
orbit cd backup schedule -o json
```

---

## cluster-profile

Manage GoCD elastic agent cluster profiles.

**Aliases:** `cp`

### cluster-profile list

List cluster profiles.

```
orbit gocd cluster-profile list [flags]
```

**Example:**

```bash
orbit gocd cluster-profile list
orbit cd cp list -o json
```

### cluster-profile get

Get cluster profile details.

```
orbit gocd cluster-profile get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Cluster profile ID |

**Example:**

```bash
orbit gocd cluster-profile get my-cluster
orbit cd cp get my-cluster -o json
```

### cluster-profile create

Create a cluster profile from a file.

```
orbit gocd cluster-profile create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with cluster profile definition |

**Example:**

```bash
orbit gocd cluster-profile create --from-file profile.yaml
orbit cd cp create --from-file profile.json
```

### cluster-profile update

Update a cluster profile from a file. Automatically fetches the current ETag.

```
orbit gocd cluster-profile update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Cluster profile ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with cluster profile configuration |

**Example:**

```bash
orbit gocd cluster-profile update my-cluster --from-file profile.yaml
orbit cd cp update my-cluster --from-file profile.json
```

### cluster-profile delete

Delete a cluster profile.

```
orbit gocd cluster-profile delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Cluster profile ID |

**Example:**

```bash
orbit gocd cluster-profile delete my-cluster
orbit cd cp delete my-cluster
```

---

## elastic-agent-profile

Manage GoCD elastic agent profiles.

**Aliases:** `eap`

### elastic-agent-profile list

List elastic agent profiles.

```
orbit gocd elastic-agent-profile list [flags]
```

**Example:**

```bash
orbit gocd elastic-agent-profile list
orbit cd eap list -o json
```

### elastic-agent-profile get

Get elastic agent profile details.

```
orbit gocd elastic-agent-profile get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

**Example:**

```bash
orbit gocd elastic-agent-profile get my-profile
orbit cd eap get my-profile -o json
```

### elastic-agent-profile create

Create an elastic agent profile from a file.

```
orbit gocd elastic-agent-profile create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with elastic agent profile definition |

**Example:**

```bash
orbit gocd elastic-agent-profile create --from-file profile.yaml
orbit cd eap create --from-file profile.json
```

### elastic-agent-profile update

Update an elastic agent profile from a file. Automatically fetches the current ETag.

```
orbit gocd elastic-agent-profile update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with elastic agent profile configuration |

**Example:**

```bash
orbit gocd elastic-agent-profile update my-profile --from-file profile.yaml
orbit cd eap update my-profile --from-file profile.json
```

### elastic-agent-profile delete

Delete an elastic agent profile.

```
orbit gocd elastic-agent-profile delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

**Example:**

```bash
orbit gocd elastic-agent-profile delete my-profile
orbit cd eap delete my-profile
```

### elastic-agent-profile usage

Show which pipelines, stages, and jobs use an elastic agent profile.

```
orbit gocd elastic-agent-profile usage <id> [flags]
```

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

**Example:**

```bash
orbit gocd elastic-agent-profile usage my-profile
orbit cd eap usage my-profile -o json
```

---

## material

Manage GoCD materials.

### material list

List all materials.

```
orbit gocd material list [flags]
```

**Example:**

```bash
orbit gocd material list
orbit cd material list -o json
```

### material get

Get material details.

```
orbit gocd material get <fingerprint>
```

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

**Example:**

```bash
orbit gocd material get abc123
orbit cd material get abc123 -o json
```

### material usage

Show pipelines using a material.

```
orbit gocd material usage <fingerprint> [flags]
```

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

**Example:**

```bash
orbit gocd material usage abc123
orbit cd material usage abc123 -o json
```

### material trigger-update

Trigger a material update check.

```
orbit gocd material trigger-update <fingerprint>
```

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

**Example:**

```bash
orbit gocd material trigger-update abc123
```

---

## artifact

Manage GoCD artifact stores.

### artifact list-store

List artifact stores.

```
orbit gocd artifact list-store [flags]
```

**Example:**

```bash
orbit gocd artifact list-store
orbit cd artifact list-store -o json
```

### artifact get-store

Get artifact store details.

```
orbit gocd artifact get-store <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Artifact store ID |

**Example:**

```bash
orbit gocd artifact get-store my-store
orbit cd artifact get-store my-store -o json
```

### artifact create-store

Create an artifact store from a file.

```
orbit gocd artifact create-store --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact store definition |

**Example:**

```bash
orbit gocd artifact create-store --from-file store.yaml
orbit cd artifact create-store --from-file store.json
```

### artifact update-store

Update an artifact store from a file. Automatically fetches the current ETag.

```
orbit gocd artifact update-store <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Artifact store ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact store configuration |

**Example:**

```bash
orbit gocd artifact update-store my-store --from-file store.yaml
orbit cd artifact update-store my-store --from-file store.json
```

### artifact delete-store

Delete an artifact store.

```
orbit gocd artifact delete-store <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Artifact store ID |

**Example:**

```bash
orbit gocd artifact delete-store my-store
```

---

## stage

Manage GoCD stages.

### stage cancel

Cancel a running stage.

```
orbit gocd stage cancel --pipeline <name> --counter <N> --stage <name> --stage-counter <N>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--counter` | Yes | Pipeline counter |
| `--stage` | Yes | Stage name |
| `--stage-counter` | Yes | Stage counter |

**Example:**

```bash
orbit gocd stage cancel --pipeline my-pipeline --counter 1 --stage my-stage --stage-counter 1
```

### stage run

Run (or re-run) a stage. Optionally specify individual jobs.

```
orbit gocd stage run --pipeline <name> --counter <N> --stage <name> [flags]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--counter` | Yes | Pipeline counter |
| `--stage` | Yes | Stage name |
| `--job` | No | Job name to run (repeatable, optional) |

**Example:**

```bash
orbit gocd stage run --pipeline my-pipeline --counter 1 --stage my-stage
orbit cd stage run --pipeline my-pipeline --counter 1 --stage my-stage --job job1 --job job2
```

---

## job

Manage GoCD jobs.

### job run

Run specific jobs in a stage.

```
orbit gocd job run --pipeline <name> --stage <name> --pipeline-counter <N> --stage-counter <N> --job <name> [flags]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--stage` | Yes | Stage name |
| `--pipeline-counter` | Yes | Pipeline counter |
| `--stage-counter` | Yes | Stage counter |
| `--job` | Yes | Job name to run (repeatable) |

**Example:**

```bash
orbit gocd job run --pipeline my-pipeline --stage my-stage --pipeline-counter 1 --stage-counter 1 --job job1 --job job2
```

---

### job log

View console output for a job run. Aliases: `logs`, `console`.

```
orbit gocd job log --pipeline <name> --stage <name> --job <name> --pipeline-counter <N> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--pipeline` | Yes | | Pipeline name. |
| `--stage` | Yes | | Stage name. |
| `--job` | Yes | | Job name. |
| `--pipeline-counter` | Yes | | Pipeline counter. |
| `--stage-counter` | No | 1 | Stage counter. |
| `--tail` | No | 0 | Show only the last N lines. |

> **404 hint:** if the job name is wrong, the CLI fetches the pipeline instance and prints the actual stage/job names on stderr before returning the error, so you can see what to pass without having to look it up separately.

**Examples:**

```bash
# Full console log
orbit cd job log --pipeline my-pipeline --stage build --job compile --pipeline-counter 42 -p myprofile

# Last 50 lines only
orbit cd job log --pipeline deploy --stage prod --job deploy-app --pipeline-counter 10 --tail 50 -p myprofile
```

---

## server

Manage GoCD server.

### server health

Show server health messages.

```
orbit gocd server health [flags]
```

**Example:**

```bash
orbit gocd server health
orbit cd server health -o json
```

### server maintenance

Show maintenance mode status.

```
orbit gocd server maintenance [flags]
```

**Example:**

```bash
orbit gocd server maintenance
orbit cd server maintenance -o json
```

### server maintenance-on

Enable maintenance mode.

```
orbit gocd server maintenance-on
```

**Example:**

```bash
orbit gocd server maintenance-on
```

### server maintenance-off

Disable maintenance mode.

```
orbit gocd server maintenance-off
```

**Example:**

```bash
orbit gocd server maintenance-off
```

---

## server-config

Manage GoCD server configuration.

**Aliases:** `sc`

### server-config site-url get

Get site URL configuration.

```
orbit gocd server-config site-url get [flags]
```

**Example:**

```bash
orbit gocd server-config site-url get
orbit cd sc site-url get -o json
```

### server-config site-url update

Update site URL configuration.

```
orbit gocd server-config site-url update [flags]
```

| Flag | Description |
|------|-------------|
| `--url` | Site URL |
| `--secure-url` | Secure site URL |

**Example:**

```bash
orbit gocd server-config site-url update --url https://gocd.example.com --secure-url https://gocd.example.com
orbit cd sc site-url update --url https://gocd.example.com
```

### server-config artifact-config get

Get server artifact configuration.

```
orbit gocd server-config artifact-config get [flags]
```

**Example:**

```bash
orbit gocd server-config artifact-config get
orbit cd sc artifact-config get -o json
```

### server-config artifact-config update

Update server artifact configuration from a file.

```
orbit gocd server-config artifact-config update --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact configuration |

**Example:**

```bash
orbit gocd server-config artifact-config update --from-file artifact-config.yaml
orbit cd sc artifact-config update --from-file artifact-config.json
```

### server-config job-timeout get

Get default job timeout.

```
orbit gocd server-config job-timeout get [flags]
```

**Example:**

```bash
orbit gocd server-config job-timeout get
orbit cd sc job-timeout get -o json
```

### server-config job-timeout update

Update default job timeout.

```
orbit gocd server-config job-timeout update --timeout <minutes>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--timeout` | Yes | Default job timeout in minutes (0 for never) |

**Example:**

```bash
orbit gocd server-config job-timeout update --timeout 60
orbit cd sc job-timeout update --timeout 0
```

### server-config mail-server get

Get mail server configuration.

```
orbit gocd server-config mail-server get [flags]
```

**Example:**

```bash
orbit gocd server-config mail-server get
orbit cd sc mail-server get -o json
```

### server-config mail-server update

Update mail server configuration from a file.

```
orbit gocd server-config mail-server update --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with mail server configuration |

**Example:**

```bash
orbit gocd server-config mail-server update --from-file mail.yaml
orbit cd sc mail-server update --from-file mail.json
```

### server-config mail-server delete

Delete mail server configuration.

```
orbit gocd server-config mail-server delete
```

**Example:**

```bash
orbit gocd server-config mail-server delete
```

---

## access-token

Manage GoCD access tokens.

**Aliases:** `token`

### access-token list

List your access tokens.

```
orbit gocd access-token list [flags]
```

**Example:**

```bash
orbit gocd access-token list
orbit cd token list -o json
```

### access-token create

Create a new access token.

```
orbit gocd access-token create --description <desc>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--description` | Yes | Description for the access token |

**Example:**

```bash
orbit gocd access-token create --description "CI token"
orbit cd token create --description "Deploy key"
```

### access-token revoke

Revoke an access token.

```
orbit gocd access-token revoke <id> [flags]
```

| Argument | Description |
|----------|-------------|
| `id` | Token ID |

| Flag | Default | Description |
|------|---------|-------------|
| `--cause` | `""` | Reason for revoking the token |

**Example:**

```bash
orbit gocd access-token revoke 42 --cause "No longer needed"
orbit cd token revoke 42
```

### access-token list-all

List all access tokens across all users (admin only).

```
orbit gocd access-token list-all [flags]
```

**Example:**

```bash
orbit gocd access-token list-all
orbit cd token list-all -o json
```

### access-token revoke-admin

Revoke any user's access token (admin only).

```
orbit gocd access-token revoke-admin <id> [flags]
```

| Argument | Description |
|----------|-------------|
| `id` | Token ID |

| Flag | Default | Description |
|------|---------|-------------|
| `--cause` | `""` | Reason for revoking the token |

**Example:**

```bash
orbit gocd access-token revoke-admin 42 --cause "Security review"
orbit cd token revoke-admin 42
```

---

## dashboard

Show the GoCD dashboard with pipeline groups, instances, and stage statuses. Compatible with both GoCD 25.x (v4 format: pipeline names as strings + objects at top-level `_embedded.pipelines`) and older versions (objects nested inside each group's `_embedded`).

```
orbit gocd dashboard [flags]
```

**Example:**

```bash
orbit gocd dashboard
orbit cd dashboard -o json
```

---

## notification-filter

Manage GoCD notification filters.

**Aliases:** `nf`

### notification-filter list

List all notification filters.

```
orbit gocd notification-filter list [flags]
```

**Example:**

```bash
orbit gocd notification-filter list
orbit cd nf list -o json
```

### notification-filter get

Get a notification filter.

```
orbit gocd notification-filter get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Filter ID |

**Example:**

```bash
orbit gocd notification-filter get 1
orbit cd nf get 1 -o json
```

### notification-filter create

Create a notification filter from a file.

```
orbit gocd notification-filter create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with filter definition |

**Example:**

```bash
orbit gocd notification-filter create --from-file filter.yaml
orbit cd nf create --from-file filter.json
```

### notification-filter delete

Delete a notification filter.

```
orbit gocd notification-filter delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Filter ID |

**Example:**

```bash
orbit gocd notification-filter delete 1
orbit cd nf delete 1
```

---

## package

Manage GoCD packages.

**Aliases:** `pkg`

### package list

List all packages.

```
orbit gocd package list [flags]
```

**Example:**

```bash
orbit gocd package list
orbit cd pkg list -o json
```

### package get

Get a package.

```
orbit gocd package get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Package ID |

**Example:**

```bash
orbit gocd package get pkg-id
orbit cd pkg get pkg-id -o json
```

### package create

Create a package from a file.

```
orbit gocd package create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with package definition |

**Example:**

```bash
orbit gocd package create --from-file package.yaml
orbit cd pkg create --from-file package.json
```

### package update

Update a package from a file. Automatically fetches the current ETag.

```
orbit gocd package update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Package ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with package definition |

**Example:**

```bash
orbit gocd package update pkg-id --from-file package.yaml
orbit cd pkg update pkg-id --from-file package.json
```

### package delete

Delete a package.

```
orbit gocd package delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Package ID |

**Example:**

```bash
orbit gocd package delete pkg-id
orbit cd pkg delete pkg-id
```

### package usage

Show pipelines using a package.

```
orbit gocd package usage <id> [flags]
```

| Argument | Description |
|----------|-------------|
| `id` | Package ID |

**Example:**

```bash
orbit gocd package usage pkg-id
orbit cd pkg usage pkg-id -o json
```

---

## package-repo

Manage GoCD package repositories.

**Aliases:** `pkg-repo`

### package-repo list

List all package repositories.

```
orbit gocd package-repo list [flags]
```

**Example:**

```bash
orbit gocd package-repo list
orbit cd pkg-repo list -o json
```

### package-repo get

Get a package repository.

```
orbit gocd package-repo get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Package repository ID |

**Example:**

```bash
orbit gocd package-repo get repo-id
orbit cd pkg-repo get repo-id -o json
```

### package-repo create

Create a package repository from a file.

```
orbit gocd package-repo create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with repository definition |

**Example:**

```bash
orbit gocd package-repo create --from-file repo.yaml
orbit cd pkg-repo create --from-file repo.json
```

### package-repo update

Update a package repository from a file. Automatically fetches the current ETag.

```
orbit gocd package-repo update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Package repository ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with repository definition |

**Example:**

```bash
orbit gocd package-repo update repo-id --from-file repo.yaml
orbit cd pkg-repo update repo-id --from-file repo.json
```

### package-repo delete

Delete a package repository.

```
orbit gocd package-repo delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Package repository ID |

**Example:**

```bash
orbit gocd package-repo delete repo-id
orbit cd pkg-repo delete repo-id
```

---

## secret-config

Manage GoCD secret configurations.

**Aliases:** `secret`

### secret-config list

List all secret configurations.

```
orbit gocd secret-config list [flags]
```

**Example:**

```bash
orbit gocd secret-config list
orbit cd secret list -o json
```

### secret-config get

Get a secret configuration.

```
orbit gocd secret-config get <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Secret configuration ID |

**Example:**

```bash
orbit gocd secret-config get my-secret
orbit cd secret get my-secret -o json
```

### secret-config create

Create a secret configuration from a file.

```
orbit gocd secret-config create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with secret config definition |

**Example:**

```bash
orbit gocd secret-config create --from-file secret.yaml
orbit cd secret create --from-file secret.json
```

### secret-config update

Update a secret configuration from a file. Automatically fetches the current ETag.

```
orbit gocd secret-config update <id> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `id` | Secret configuration ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with secret config definition |

**Example:**

```bash
orbit gocd secret-config update my-secret --from-file secret.yaml
orbit cd secret update my-secret --from-file secret.json
```

### secret-config delete

Delete a secret configuration.

```
orbit gocd secret-config delete <id>
```

| Argument | Description |
|----------|-------------|
| `id` | Secret configuration ID |

**Example:**

```bash
orbit gocd secret-config delete my-secret
orbit cd secret delete my-secret
```

---

## template

Manage GoCD pipeline templates.

**Aliases:** `tmpl`

### template list

List all pipeline templates.

```
orbit gocd template list [flags]
```

**Example:**

```bash
orbit gocd template list
orbit cd tmpl list -o json
```

### template get

Get a pipeline template.

```
orbit gocd template get <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Template name |

**Example:**

```bash
orbit gocd template get my-template
orbit cd tmpl get my-template -o json
```

### template create

Create a pipeline template from a file.

```
orbit gocd template create --from-file <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with template definition |

**Example:**

```bash
orbit gocd template create --from-file template.yaml
orbit cd tmpl create --from-file template.json
```

### template update

Update a pipeline template from a file. Automatically fetches the current ETag.

```
orbit gocd template update <name> --from-file <file>
```

| Argument | Description |
|----------|-------------|
| `name` | Template name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with template definition |

**Example:**

```bash
orbit gocd template update my-template --from-file template.yaml
orbit cd tmpl update my-template --from-file template.json
```

### template delete

Delete a pipeline template.

```
orbit gocd template delete <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Template name |

**Example:**

```bash
orbit gocd template delete my-template
orbit cd tmpl delete my-template
```

---

## version

Show GoCD server version.

```
orbit gocd version [flags]
```

**Example:**

```bash
orbit gocd version
orbit cd version -o json
```

---

## encrypt

Encrypt a value using GoCD's cipher.

```
orbit gocd encrypt <value>
```

| Argument | Description |
|----------|-------------|
| `value` | Plain-text value to encrypt |

**Example:**

```bash
orbit gocd encrypt my-secret-password
orbit cd encrypt my-database-password -o json
```
