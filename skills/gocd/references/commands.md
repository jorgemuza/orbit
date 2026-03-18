# GoCD Command Reference

Complete reference for all `orbit gocd` (alias: `cd`) commands.

## Pipeline Commands

### `orbit gocd pipeline list`
List all pipelines grouped by pipeline group.

| Flag | Default | Description |
|------|---------|-------------|
| `-o` | `table` | Output format: table, json, yaml |

### `orbit gocd pipeline status <name>`
Show pipeline status including paused, schedulable, and locked state.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

### `orbit gocd pipeline history <name>`
Show pipeline run history.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `10` | Maximum number of history entries |

### `orbit gocd pipeline get <name>`
Get a pipeline instance.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--counter` | `0` | Pipeline counter (0 = latest) |

### `orbit gocd pipeline trigger <name>`
Schedule a pipeline run.

### `orbit gocd pipeline pause <name>`
Pause a pipeline.

| Flag | Default | Description |
|------|---------|-------------|
| `--reason` | `""` | Reason for pausing |

### `orbit gocd pipeline unpause <name>`
Unpause a pipeline.

### `orbit gocd pipeline create --group <group> --from-file <file>`
Create a pipeline from a JSON or YAML file.

| Flag | Required | Description |
|------|----------|-------------|
| `--group` | Yes | Pipeline group name |
| `--from-file` | Yes | Path to JSON or YAML file with pipeline definition |

### `orbit gocd pipeline update <name> --from-file <file>`
Update a pipeline configuration. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline configuration |

### `orbit gocd pipeline delete <name>`
Delete a pipeline.

### `orbit gocd pipeline comment <name> --counter <N> --message <msg>`
Add a comment to a pipeline instance.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Required | Description |
|------|----------|-------------|
| `--counter` | Yes | Pipeline counter |
| `--message` | Yes | Comment message |

### `orbit gocd pipeline export <name>`
Export pipeline configuration.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline name |

| Flag | Default | Description |
|------|---------|-------------|
| `--plugin-id` | `""` | Config repo plugin ID (optional) |

## Pipeline Group Commands

Aliases: `pipeline-group`, `pg`

### `orbit gocd pipeline-group list`
List all pipeline groups with name and pipeline count.

### `orbit gocd pipeline-group get <name>`
Get pipeline group details.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline group name |

### `orbit gocd pipeline-group create --from-file <file>`
Create a pipeline group from a JSON or YAML file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline group definition |

### `orbit gocd pipeline-group update <name> --from-file <file>`
Update a pipeline group. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `name` | Pipeline group name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with pipeline group configuration |

### `orbit gocd pipeline-group delete <name>`
Delete a pipeline group.

## Agent Commands

### `orbit gocd agent list`
List all agents with UUID, hostname, OS, config state, agent state, and build state.

### `orbit gocd agent get <uuid>`
View detailed agent information.

### `orbit gocd agent enable <uuid>`
Enable a disabled agent. Uses ETag-based optimistic locking.

### `orbit gocd agent disable <uuid>`
Disable an agent. Uses ETag-based optimistic locking.

### `orbit gocd agent delete <uuid>`
Delete an agent permanently.

### `orbit gocd agent kill-task <uuid>`
Kill running tasks on an agent.

### `orbit gocd agent update <uuid> --from-file <file>`
Update an agent from a JSON or YAML file.

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with agent update data |

### `orbit gocd agent job-history <uuid>`
Show job run history for an agent. Displays pipeline, counter, stage, job, state, and result.

| Argument | Description |
|----------|-------------|
| `uuid` | Agent UUID |

## Environment Commands

Aliases: `environment`, `env`

### `orbit gocd environment list`
List all environments with their pipelines and agents.

### `orbit gocd environment get <name>`
View detailed environment information including pipelines, agents, and environment variables.

### `orbit gocd environment create <name>`
Create a new environment.

| Flag | Description |
|------|-------------|
| `--pipeline` | Pipeline to add (repeatable) |

### `orbit gocd environment delete <name>`
Delete an environment.

### `orbit gocd environment update <name> --from-file <file>`
Replace an environment definition from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with environment definition |

### `orbit gocd environment patch <name>`
Patch an environment by adding or removing pipelines and agents.

| Argument | Description |
|----------|-------------|
| `name` | Environment name |

| Flag | Description |
|------|-------------|
| `--add-pipeline` | Pipeline to add (repeatable) |
| `--remove-pipeline` | Pipeline to remove (repeatable) |
| `--add-agent` | Agent to add (repeatable) |
| `--remove-agent` | Agent to remove (repeatable) |

## Config Repo Commands

Aliases: `config-repo`, `configrepo`, `cr`

### `orbit gocd config-repo list`
List all config repos with ID, plugin, and material type.

### `orbit gocd config-repo get <id>`
View config repo details.

### `orbit gocd config-repo status <id>`
Check config repo sync status (in progress, state, errors).

### `orbit gocd config-repo trigger <id>`
Trigger a config repo update/re-sync.

### `orbit gocd config-repo create --from-file <file>`
Create a config repository from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with config repo definition |

### `orbit gocd config-repo update <id> --from-file <file>`
Update a config repository from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with config repo definition |

### `orbit gocd config-repo delete <id>`
Delete a config repository.

### `orbit gocd config-repo definitions <id>`
Get the pipelines and environments defined by a config repo.

| Argument | Description |
|----------|-------------|
| `id` | Config repo ID |

### `orbit gocd config-repo preflight --plugin-id <id> --from-file <file>`
Run a preflight check on config repo content before applying it.

| Flag | Required | Description |
|------|----------|-------------|
| `--plugin-id` | Yes | Config repo plugin ID |
| `--from-file` | Yes | Path to JSON or YAML file with content to check |
| `--repo-id` | No | Config repo ID (optional, for existing repos) |

## User Commands

### `orbit gocd user list`
List all users with login name, display name, email, enabled state, and admin flag.

### `orbit gocd user get <login>`
View detailed user information.

| Argument | Description |
|----------|-------------|
| `login` | User login name |

### `orbit gocd user create --from-file <file>`
Create a user from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with user definition |

### `orbit gocd user update <login> --from-file <file>`
Update a user from a file.

| Argument | Description |
|----------|-------------|
| `login` | User login name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with user configuration |

### `orbit gocd user delete <login>`
Delete a user.

### `orbit gocd user delete-bulk`
Delete multiple users at once.

| Flag | Description |
|------|-------------|
| `--user` | User login name to delete (repeatable) |
| `--from-file` | Path to JSON or YAML file with bulk delete request |

## Plugin Commands

### `orbit gocd plugin list`
List installed plugins with ID, bundled status, and state.

### `orbit gocd plugin get <id>`
View detailed plugin information.

| Argument | Description |
|----------|-------------|
| `id` | Plugin ID |

### `orbit gocd plugin get-settings <plugin-id>`
Get plugin settings.

| Argument | Description |
|----------|-------------|
| `plugin-id` | Plugin ID |

### `orbit gocd plugin create-settings --from-file <file>`
Create plugin settings from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with plugin settings |

### `orbit gocd plugin update-settings <plugin-id> --from-file <file>`
Update plugin settings from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `plugin-id` | Plugin ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with plugin settings |

## Role Commands

### `orbit gocd role list`
List all security roles with name and type.

### `orbit gocd role get <name>`
Get role details.

| Argument | Description |
|----------|-------------|
| `name` | Role name |

### `orbit gocd role create --from-file <file>`
Create a role from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with role definition |

### `orbit gocd role update <name> --from-file <file>`
Update a role from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `name` | Role name |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with role configuration |

### `orbit gocd role delete <name>`
Delete a role.

## Authorization Commands

Aliases: `authorization`, `auth-config`

### `orbit gocd authorization list`
List authorization configurations with ID and plugin ID.

### `orbit gocd authorization get <id>`
Get authorization configuration details.

| Argument | Description |
|----------|-------------|
| `id` | Authorization configuration ID |

### `orbit gocd authorization create --from-file <file>`
Create an authorization configuration from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with authorization configuration |

### `orbit gocd authorization update <id> --from-file <file>`
Update an authorization configuration from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `id` | Authorization configuration ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with authorization configuration |

### `orbit gocd authorization delete <id>`
Delete an authorization configuration.

## Backup Commands

### `orbit gocd backup get-config`
Get backup configuration (schedule, email settings, post-backup script).

### `orbit gocd backup create-config --from-file <file>`
Create backup configuration from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with backup configuration |

### `orbit gocd backup delete-config`
Delete backup configuration.

### `orbit gocd backup schedule`
Schedule a server backup. Returns status and path.

## Cluster Profile Commands

Aliases: `cluster-profile`, `cp`

### `orbit gocd cluster-profile list`
List cluster profiles with ID and plugin ID.

### `orbit gocd cluster-profile get <id>`
Get cluster profile details.

| Argument | Description |
|----------|-------------|
| `id` | Cluster profile ID |

### `orbit gocd cluster-profile create --from-file <file>`
Create a cluster profile from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with cluster profile definition |

### `orbit gocd cluster-profile update <id> --from-file <file>`
Update a cluster profile from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `id` | Cluster profile ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with cluster profile configuration |

### `orbit gocd cluster-profile delete <id>`
Delete a cluster profile.

## Elastic Agent Profile Commands

Aliases: `elastic-agent-profile`, `eap`

### `orbit gocd elastic-agent-profile list`
List elastic agent profiles with ID and cluster profile ID.

### `orbit gocd elastic-agent-profile get <id>`
Get elastic agent profile details.

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

### `orbit gocd elastic-agent-profile create --from-file <file>`
Create an elastic agent profile from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with elastic agent profile definition |

### `orbit gocd elastic-agent-profile update <id> --from-file <file>`
Update an elastic agent profile from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with elastic agent profile configuration |

### `orbit gocd elastic-agent-profile delete <id>`
Delete an elastic agent profile.

### `orbit gocd elastic-agent-profile usage <id>`
Show which pipelines, stages, and jobs use an elastic agent profile.

| Argument | Description |
|----------|-------------|
| `id` | Elastic agent profile ID |

## Material Commands

### `orbit gocd material list`
List all materials with type and fingerprint.

### `orbit gocd material get <fingerprint>`
Get material details including type, fingerprint, and attributes.

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

### `orbit gocd material usage <fingerprint>`
Show pipelines using a material.

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

### `orbit gocd material trigger-update <fingerprint>`
Trigger a material update check.

| Argument | Description |
|----------|-------------|
| `fingerprint` | Material fingerprint |

## Artifact Store Commands

### `orbit gocd artifact list-store`
List artifact stores with ID and plugin ID.

### `orbit gocd artifact get-store <id>`
Get artifact store details.

| Argument | Description |
|----------|-------------|
| `id` | Artifact store ID |

### `orbit gocd artifact create-store --from-file <file>`
Create an artifact store from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact store definition |

### `orbit gocd artifact update-store <id> --from-file <file>`
Update an artifact store from a file. Automatically fetches the current ETag.

| Argument | Description |
|----------|-------------|
| `id` | Artifact store ID |

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact store configuration |

### `orbit gocd artifact delete-store <id>`
Delete an artifact store.

## Stage Commands

### `orbit gocd stage cancel`
Cancel a running stage.

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--counter` | Yes | Pipeline counter |
| `--stage` | Yes | Stage name |
| `--stage-counter` | Yes | Stage counter |

### `orbit gocd stage run`
Run (or re-run) a stage. Optionally specify individual jobs.

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--counter` | Yes | Pipeline counter |
| `--stage` | Yes | Stage name |
| `--job` | No | Job name to run (repeatable) |

## Job Commands

### `orbit gocd job run`
Run specific jobs in a stage.

| Flag | Required | Description |
|------|----------|-------------|
| `--pipeline` | Yes | Pipeline name |
| `--stage` | Yes | Stage name |
| `--pipeline-counter` | Yes | Pipeline counter |
| `--stage-counter` | Yes | Stage counter |
| `--job` | Yes | Job name to run (repeatable) |

## Server Commands

### `orbit gocd server health`
Show server health messages with level, message, detail, and time.

### `orbit gocd server maintenance`
Show current maintenance mode status.

### `orbit gocd server maintenance-on`
Enable maintenance mode.

### `orbit gocd server maintenance-off`
Disable maintenance mode.

## Server Config Commands

Aliases: `server-config`, `sc`

### `orbit gocd server-config site-url get`
Get site URL configuration (site URL and secure site URL).

### `orbit gocd server-config site-url update`
Update site URL configuration.

| Flag | Description |
|------|-------------|
| `--url` | Site URL |
| `--secure-url` | Secure site URL |

### `orbit gocd server-config artifact-config get`
Get server artifact configuration (artifacts dir, purge settings).

### `orbit gocd server-config artifact-config update --from-file <file>`
Update server artifact configuration from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with artifact configuration |

### `orbit gocd server-config job-timeout get`
Get default job timeout.

### `orbit gocd server-config job-timeout update --timeout <minutes>`
Update default job timeout.

| Flag | Required | Description |
|------|----------|-------------|
| `--timeout` | Yes | Default job timeout in minutes (0 for never) |

### `orbit gocd server-config mail-server get`
Get mail server configuration.

### `orbit gocd server-config mail-server update --from-file <file>`
Update mail server configuration from a file.

| Flag | Required | Description |
|------|----------|-------------|
| `--from-file` | Yes | Path to JSON or YAML file with mail server configuration |

### `orbit gocd server-config mail-server delete`
Delete mail server configuration.

## Encryption

### `orbit gocd encrypt <value>`
Encrypt a plain-text value using GoCD's cipher. Returns the encrypted value for use in GoCD configuration files.

## Global Flags

| Flag | Description |
|------|-------------|
| `--service` | GoCD service name (only needed with multiple GoCD services in profile) |
| `-p, --profile` | Profile to use |
| `-o, --output` | Output format: table, json, yaml |
| `--config` | Config file path (default: `~/.config/orbit/config.yaml`) |
