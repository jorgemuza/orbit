---
name: gocd
description: "Manage GoCD pipelines, pipeline groups, agents, environments, config repos, server administration, users, roles, authorization configs, plugins, backups, materials, artifact stores, elastic agent profiles, cluster profiles, stages, jobs, and server configuration using the orbit CLI. Use this skill whenever the user asks about GoCD pipelines, agents, environments, config repos, server health, maintenance mode, CI/CD operations, pipeline groups, users, roles, auth configs, plugins, backups, materials, artifact stores, elastic agents, stages, or jobs on GoCD. Trigger on phrases like 'list pipelines', 'create pipeline', 'delete pipeline', 'pipeline status', 'trigger a build', 'pause pipeline', 'list agents', 'enable agent', 'disable agent', 'kill running tasks', 'agent job history', 'list environments', 'create environment', 'patch environment', 'config repo status', 'create config repo', 'preflight check', 'server health', 'maintenance mode', 'encrypt a value', 'list users', 'create user', 'list roles', 'auth config', 'list plugins', 'schedule backup', 'list materials', 'artifact store', 'elastic agent profile', 'cluster profile', 'cancel stage', 'run stage', 'run job', 'site url', 'job timeout', 'mail server config', 'pipeline group', or any GoCD-related task — even casual references like 'what pipelines are running', 'is the agent idle', 'check the build', 'schedule a run', 'put server in maintenance', 'check config repo sync', 'who has access', 'what plugins are installed', or 'GoCD status'. The orbit CLI alias is `cd`."
---

# GoCD with orbit CLI

Comprehensive GoCD management through the `orbit` CLI. Covers pipelines, pipeline groups, agents, environments, config repos, server administration, users, roles, authorization, plugins, backups, materials, artifact stores, elastic agent profiles, cluster profiles, stages, jobs, and server configuration. Works with self-hosted GoCD instances via the REST API, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` binary built and accessible
2. A profile with a `gocd` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials (API token or basic auth) — can be stored in 1Password with `op://` prefix
4. `base_url` is required (GoCD is always self-hosted)

## Quick Reference

All commands follow the pattern: `orbit -p <profile> gocd <command> [flags]`

Alias: `orbit -p <profile> cd <command> [flags]`

All commands support `-o json` and `-o yaml` for structured output. For full command details and all flags, see `references/commands.md`.

## Command Groups

| Group | Alias | Description |
|-------|-------|-------------|
| `pipeline` | | Pipeline CRUD, status, history, trigger, pause/unpause, comment, export |
| `pipeline-group` | `pg` | Pipeline group CRUD |
| `agent` | | Agent list, get, enable, disable, delete, kill-task, update, job-history |
| `environment` | `env` | Environment CRUD, patch |
| `config-repo` | `cr`, `configrepo` | Config repo CRUD, status, trigger, definitions, preflight |
| `server` | | Server health, maintenance mode |
| `server-config` | `sc` | Site URL, artifact config, job timeout, mail server |
| `user` | | User CRUD, bulk delete |
| `role` | | Security role CRUD |
| `authorization` | `auth-config` | Authorization config CRUD |
| `plugin` | | Plugin info, settings CRUD |
| `backup` | | Backup config CRUD, schedule |
| `material` | | Material list, get, usage, trigger-update |
| `artifact` | | Artifact store CRUD |
| `cluster-profile` | `cp` | Elastic agent cluster profile CRUD |
| `elastic-agent-profile` | `eap` | Elastic agent profile CRUD, usage |
| `stage` | | Stage cancel, run |
| `job` | | Job run |
| `encrypt` | | Encrypt values |

## Core Workflows

### Pipeline CRUD

```bash
# List all pipelines by group
orbit -p myprofile gocd pipeline list

# Create a pipeline from file
orbit -p myprofile cd pipeline create --group my-group --from-file pipeline.yaml

# Update a pipeline configuration
orbit -p myprofile cd pipeline update my-pipeline --from-file pipeline.yaml

# Delete a pipeline
orbit -p myprofile cd pipeline delete my-pipeline

# Check pipeline status
orbit -p myprofile cd pipeline status my-pipeline

# View pipeline run history
orbit -p myprofile gocd pipeline history my-pipeline --limit 5

# Get a specific pipeline instance
orbit -p myprofile gocd pipeline get my-pipeline --counter 42

# Trigger a pipeline run
orbit -p myprofile cd pipeline trigger my-pipeline

# Pause/unpause a pipeline
orbit -p myprofile gocd pipeline pause my-pipeline --reason "Maintenance window"
orbit -p myprofile gocd pipeline unpause my-pipeline

# Comment on a pipeline instance
orbit -p myprofile cd pipeline comment my-pipeline --counter 42 --message "Deployed to staging"

# Export pipeline configuration
orbit -p myprofile cd pipeline export my-pipeline --plugin-id yaml.config.plugin
```

### Pipeline Group Management

```bash
# List pipeline groups
orbit -p myprofile gocd pipeline-group list
orbit -p myprofile cd pg list

# Get pipeline group details
orbit -p myprofile cd pg get my-group

# Create/update/delete pipeline groups
orbit -p myprofile cd pg create --from-file group.yaml
orbit -p myprofile cd pg update my-group --from-file group.yaml
orbit -p myprofile cd pg delete my-group
```

### Agent Management

```bash
# List all agents
orbit -p myprofile gocd agent list

# View agent details
orbit -p myprofile cd agent get <uuid>

# Enable/disable agents
orbit -p myprofile gocd agent enable <uuid>
orbit -p myprofile gocd agent disable <uuid>

# Update agent configuration
orbit -p myprofile cd agent update <uuid> --from-file agent.yaml

# View agent job history
orbit -p myprofile cd agent job-history <uuid>

# Kill running tasks on an agent
orbit -p myprofile cd agent kill-task <uuid>

# Delete an agent
orbit -p myprofile gocd agent delete <uuid>
```

### Environment Management

```bash
# List environments
orbit -p myprofile gocd environment list

# View environment details
orbit -p myprofile cd env get production

# Create environment with pipelines
orbit -p myprofile gocd environment create staging --pipeline my-app --pipeline my-api

# Update environment from file
orbit -p myprofile cd env update staging --from-file env.yaml

# Patch environment (add/remove pipelines and agents)
orbit -p myprofile cd env patch staging --add-pipeline my-app --remove-pipeline old-app --add-agent agent1

# Delete environment
orbit -p myprofile gocd environment delete staging
```

### Config Repo Operations

```bash
# List config repos
orbit -p myprofile gocd cr list

# Get config repo details
orbit -p myprofile cd cr get my-repo

# Create/update/delete config repos
orbit -p myprofile cd cr create --from-file repo.yaml
orbit -p myprofile cd cr update my-repo --from-file repo.yaml
orbit -p myprofile cd cr delete my-repo

# Check sync status
orbit -p myprofile cd cr status my-repo

# Trigger config repo update
orbit -p myprofile cd cr trigger my-repo

# View definitions from a config repo
orbit -p myprofile cd cr definitions my-repo

# Run preflight check
orbit -p myprofile cd cr preflight --plugin-id yaml.config.plugin --from-file pipeline.yaml
```

### User Management

```bash
# List users
orbit -p myprofile gocd user list

# Get user details
orbit -p myprofile cd user get admin

# Create/update/delete users
orbit -p myprofile cd user create --from-file user.yaml
orbit -p myprofile cd user update admin --from-file user.yaml
orbit -p myprofile cd user delete admin

# Bulk delete users
orbit -p myprofile cd user delete-bulk --user user1 --user user2
```

### Security: Roles & Authorization

```bash
# List roles
orbit -p myprofile gocd role list

# CRUD roles
orbit -p myprofile cd role get my-role
orbit -p myprofile cd role create --from-file role.yaml
orbit -p myprofile cd role update my-role --from-file role.yaml
orbit -p myprofile cd role delete my-role

# List authorization configs
orbit -p myprofile cd auth-config list

# CRUD authorization configs
orbit -p myprofile cd auth-config get my-auth
orbit -p myprofile cd auth-config create --from-file auth.yaml
orbit -p myprofile cd auth-config update my-auth --from-file auth.yaml
orbit -p myprofile cd auth-config delete my-auth
```

### Plugin Management

```bash
# List plugins
orbit -p myprofile gocd plugin list

# Get plugin info
orbit -p myprofile cd plugin get my-plugin-id

# Plugin settings CRUD
orbit -p myprofile cd plugin get-settings my-plugin-id
orbit -p myprofile cd plugin create-settings --from-file settings.yaml
orbit -p myprofile cd plugin update-settings my-plugin-id --from-file settings.yaml
```

### Backup Management

```bash
# View backup configuration
orbit -p myprofile gocd backup get-config

# Create/delete backup config
orbit -p myprofile cd backup create-config --from-file backup.yaml
orbit -p myprofile cd backup delete-config

# Schedule a backup
orbit -p myprofile cd backup schedule
```

### Material Operations

```bash
# List materials
orbit -p myprofile gocd material list

# Get material details
orbit -p myprofile cd material get <fingerprint>

# View material usage (which pipelines use it)
orbit -p myprofile cd material usage <fingerprint>

# Trigger material update check
orbit -p myprofile cd material trigger-update <fingerprint>
```

### Artifact Store Management

```bash
# List artifact stores
orbit -p myprofile gocd artifact list-store

# CRUD artifact stores
orbit -p myprofile cd artifact get-store my-store
orbit -p myprofile cd artifact create-store --from-file store.yaml
orbit -p myprofile cd artifact update-store my-store --from-file store.yaml
orbit -p myprofile cd artifact delete-store my-store
```

### Elastic Agent Infrastructure

```bash
# Cluster profiles
orbit -p myprofile cd cp list
orbit -p myprofile cd cp get my-cluster
orbit -p myprofile cd cp create --from-file cluster.yaml
orbit -p myprofile cd cp update my-cluster --from-file cluster.yaml
orbit -p myprofile cd cp delete my-cluster

# Elastic agent profiles
orbit -p myprofile cd eap list
orbit -p myprofile cd eap get my-profile
orbit -p myprofile cd eap create --from-file profile.yaml
orbit -p myprofile cd eap update my-profile --from-file profile.yaml
orbit -p myprofile cd eap delete my-profile
orbit -p myprofile cd eap usage my-profile
```

### Stage & Job Operations

```bash
# Cancel a running stage
orbit -p myprofile cd stage cancel --pipeline my-pipeline --counter 5 --stage build --stage-counter 1

# Run/re-run a stage
orbit -p myprofile cd stage run --pipeline my-pipeline --counter 5 --stage build

# Run specific jobs in a stage
orbit -p myprofile cd job run --pipeline my-pipeline --stage build --pipeline-counter 5 --stage-counter 1 --job unit-test --job integration-test
```

### Server Administration

```bash
# Server health
orbit -p myprofile gocd server health

# Maintenance mode
orbit -p myprofile cd server maintenance
orbit -p myprofile cd server maintenance-on
orbit -p myprofile cd server maintenance-off

# Site URLs
orbit -p myprofile cd sc site-url get
orbit -p myprofile cd sc site-url update --url https://gocd.example.com --secure-url https://gocd.example.com

# Artifact config
orbit -p myprofile cd sc artifact-config get
orbit -p myprofile cd sc artifact-config update --from-file artifact.yaml

# Job timeout
orbit -p myprofile cd sc job-timeout get
orbit -p myprofile cd sc job-timeout update --timeout 30

# Mail server
orbit -p myprofile cd sc mail-server get
orbit -p myprofile cd sc mail-server update --from-file mail.yaml
orbit -p myprofile cd sc mail-server delete

# Encryption
orbit -p myprofile gocd encrypt my-secret-password
```

## Common Patterns

### Create/Update Resources from File

Most CRUD commands accept `--from-file` for YAML or JSON input:
```bash
orbit -p myprofile cd pipeline create --group my-group --from-file pipeline.yaml
orbit -p myprofile cd cr create --from-file config-repo.json
orbit -p myprofile cd role create --from-file role.yaml
```

### ETag-Based Updates

For resources that require optimistic locking (pipelines, pipeline groups, environments, config repos, roles, auth configs, plugins, cluster profiles, elastic agent profiles, artifact stores), the CLI handles ETag retrieval automatically — just provide the `--from-file` with updated config.

### Debug a Failing Config Repo
```bash
orbit -p myprofile cd cr status my-config-repo
orbit -p myprofile cd cr definitions my-config-repo
orbit -p myprofile cd cr trigger my-config-repo
```

### Prepare for Maintenance
```bash
orbit -p myprofile cd agent list
orbit -p myprofile cd server maintenance-on
# ... perform maintenance ...
orbit -p myprofile cd server maintenance-off
```

## Important Notes

- GoCD API uses versioned Accept headers — handled automatically by the CLI.
- Update operations use ETag-based optimistic locking — the CLI handles this transparently.
- The `--service` flag is only needed when a profile has multiple GoCD services configured.
- `base_url` is required in service configuration (GoCD is always self-hosted).
- All destructive operations execute immediately without confirmation prompts.
- Output formats: `table` (default), `json`, `yaml`.
