---
name: jira
description: "Interact with Jira using the orbit CLI to create, list, view, edit, and transition issues, manage sprints and epics, manage custom fields and screen configurations, list statuses and issue types, and write properly formatted descriptions using Jira wiki markup. Use this skill whenever the user asks about Jira tasks, tickets, issues, sprints, epics, or needs to manage project work items using orbit. Also trigger when the user says things like 'create a ticket', 'create epics', 'move this to done', 'assign the issue', 'update the description', 'format for Jira', 'create a custom field', 'add field to screen', 'list statuses', 'configure Jira', or any Jira-related workflow — even casual references like 'update Jira', 'what tickets are in this sprint', 'add a comment to PROJ-123', or 'set up AI tracking fields'. Trigger especially when descriptions need proper formatting (headings, bullets, tables, links) since Jira Server uses wiki markup, not markdown."
---

# Jira with orbit CLI

Manage Jira issues, epics, sprints, boards, projects, releases, custom fields, screens, statuses, and issue types using the `orbit` CLI. Supports both Jira Cloud (REST API v3) and Jira Server/Data Center (REST API v2) with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` binary built and accessible
2. A profile with a `jira-cloud` or `jira-onprem` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials — API token for Cloud (basic auth with email + token), PAT for Server
4. Credentials can be stored in 1Password with `op://` prefix for automatic resolution

## Quick Reference

All commands follow the pattern: `orbit -p <profile> jira <resource> <action> [flags]`

The `-o` flag controls output format: `table` (default), `json`, `yaml`.

For full command details and flags, see `references/commands.md`.
For Jira wiki markup formatting (Server only), see `references/wiki-markup.md`.

## Core Workflows

### Creating Issues

```bash
# Simple story
orbit -p myprofile jira issue create --project PROJ --type Story --summary "Add login page"

# Bug with priority and assignment
orbit -p myprofile jira issue create --project PROJ --type Bug \
  --summary "Fix timeout" --priority High --assignee john.doe

# Sub-task under a parent
orbit -p myprofile jira issue create --project PROJ --type Sub-task \
  --parent PROJ-123 --summary "Add validation"
```

### Creating Epics

Epics require the "Epic Name" custom field. The `orbit` CLI auto-sets this from the summary when `--type Epic` is used. Use `--epic-name` to override.

```bash
# Create an epic (Epic Name auto-set from summary)
orbit -p myprofile jira issue create --type Epic --project PROJ \
  --summary "Q1 Auth Revamp" --priority Highest

# Create epic with Parent Link to a Capability/Initiative
orbit -p myprofile jira issue create --type Epic --project PROJ \
  --summary "Okta Authentication Foundation" \
  --field "customfield_27521=PRT-4378" \
  --priority Highest
```

The `--field` flag sets arbitrary custom fields as `key=value`. Use `jira field list` to discover field IDs.

### Editing Issues with Formatted Descriptions

Jira Server uses **wiki markup**, not Markdown. Jira Cloud uses ADF (handled automatically by orbit for plain text). Always format `--body` values using wiki markup syntax for Server instances.

```bash
orbit -p myprofile jira issue edit PRT-123 --body "h2. Value Statement

The platform provides real authentication via Okta.

h2. User Stories

* *Story 1:* As a platform admin, I want all API requests to require a valid Okta JWT.
** Okta custom AS configured with epsilon_claims claim
** AUTH_BYPASS removed from all Terraform task definitions"
```

### Editing Custom Field Values

The `-F` flag supports multiple value formats. The correct format depends on the Jira field type:

```bash
# Single-select field — use value: shorthand
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10397=value:Yes

# Multi-select field — MUST use JSON array of objects (value: shorthand will return 400)
orbit -p myprofile jira issue edit PROJ-123 -F 'customfield_10398=[{"value":"Developer"},{"value":"QA"}]'

# Multi-select with a single value still needs array format
orbit -p myprofile jira issue edit PROJ-123 -F 'customfield_10400=[{"value":"Confirmed"}]'

# Number field
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10399=3

# Plain string field
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10010=my-value
```

**Important**: Multi-select fields (checkboxes, multi-select dropdowns) require `[{"value":"X"}]` array format. Using `value:X` on a multi-select field returns a 400 error.

### Transitioning Issues

```bash
# Move to In Progress
orbit -p myprofile jira issue move PROJ-123 "In Progress"

# Close with comment and resolution
orbit -p myprofile jira issue move PROJ-123 Done --comment "Fixed in v2.1" --resolution Fixed
```

### Comments with @mentions

Use `[~accountId:ID]` to @mention users in comments (Cloud). Jira renders these as clickable user links. For Server, use `[~username]`.

```bash
# Add a comment mentioning the assignee
orbit -p myprofile jira issue comment PROJ-123 -b "Hey [~accountId:5b10ac8d82e05b22cc7d4ef5], this needs rework"

# Look up a user's accountId
orbit -p myprofile jira user search "Jorge"
```

### Searching Issues

```bash
# Filter by project and type
orbit -p myprofile jira issue list --project PROJ --type Epic

# Filter by assignee and status
orbit -p myprofile jira issue list --assignee me --status "In Progress"

# Raw JQL query
orbit -p myprofile jira issue list --jql "project = PROJ AND sprint in openSprints()"

# JSON output for processing
orbit -p myprofile jira issue list --project PROJ -o json
```

## Field Management (Cloud Only)

Manage custom fields, their contexts, and options programmatically.

### Listing Fields

```bash
# List all fields
orbit -p myprofile jira field list

# List only custom fields
orbit -p myprofile jira field list --custom

# Filter by name or ID
orbit -p myprofile jira field list --filter "AI"
```

### Creating Custom Fields

```bash
# Create a select field (shorthand types: select, multiselect, number, checkbox, text, textarea)
orbit -p myprofile jira field create --name "AI Assisted" --type select \
  --description "Was AI used on this ticket?"

# Create a number field
orbit -p myprofile jira field create --name "AI Prompt Iterations" --type number \
  --description "How many prompt cycles to get working output"

# Create a checkbox field
orbit -p myprofile jira field create --name "Human Review Confirmed" --type checkbox \
  --description "Engineer confirms AI output was reviewed"
```

### Managing Field Contexts and Options

```bash
# List field contexts
orbit -p myprofile jira field context-list customfield_10397

# List existing options
orbit -p myprofile jira field option-list customfield_10397 10817

# Add options to a select/multiselect field
orbit -p myprofile jira field option-add customfield_10397 10817 \
  --values "Yes,No,Partial"
```

## Screen Management

Control which fields appear on issue create/edit/view forms by managing screen tabs and field assignments.

### Listing Screens and Tabs

```bash
# List all screens
orbit -p myprofile jira screen list

# Filter screens by name
orbit -p myprofile jira screen list --filter "PYMT"

# List tabs on a screen
orbit -p myprofile jira screen tab-list 10089

# List fields on a screen tab
orbit -p myprofile jira screen field-list 10089 10189
```

### Creating Tabs

```bash
# Create a new tab to group related fields
orbit -p myprofile jira screen tab-create 10089 "AI Workflow"
```

### Adding and Removing Fields

```bash
# Add fields to a screen tab
orbit -p myprofile jira screen field-add 10089 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399"

# Remove fields from a screen tab
orbit -p myprofile jira screen field-remove 10089 10189 \
  --fields "customfield_10397,customfield_10398"
```

### Moving Fields Between Tabs

```bash
# Move fields from General tab to AI Workflow tab
orbit -p myprofile jira screen field-move 10089 10189 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399"
```

## Status and Issue Type Management

### Listing Workflow Statuses

```bash
# List all statuses with their categories
orbit -p myprofile jira status list

# Filter with grep
orbit -p myprofile jira status list | grep -i "review"
```

### Listing Issue Types

```bash
orbit -p myprofile jira issuetype-list
```

## Deleting Issues

Delete issues with `issue delete`. Use `--cascade` to also delete all subtasks.

**Important:** There is no `--confirm` or `--force` flag. The command deletes immediately.

```bash
# Delete a single issue
orbit -p myprofile jira issue delete PROJ-123

# Delete an issue and all its subtasks
orbit -p myprofile jira issue delete PROJ-100 --cascade
```

**Deleting a hierarchy (parent → children → grandchildren):** Delete bottom-up to avoid orphan errors, or use `--cascade` at each level that has subtasks.

```bash
# Example: delete Story (with sub-tasks) → Epic → Capability
orbit -p myprofile jira issue delete PROJ-50 --cascade   # Story + sub-tasks
orbit -p myprofile jira issue delete PROJ-40              # Epic (now childless)
orbit -p myprofile jira issue delete PRT-200              # Capability (now childless)
```

**Tip:** `--cascade` only deletes Jira sub-tasks (child issues of type Sub-task). It does NOT delete Stories under an Epic or Epics under a Capability. For full hierarchy deletion, query children with JQL first, then delete bottom-up.

```bash
# Find all children of an issue
orbit -p myprofile jira issue list --jql '"Parent" = PROJ-40' -o json
```

## Important Notes

- **Cloud vs Server** — Use service type `jira-cloud` for Atlassian Cloud (uses API v3 with ADF for descriptions). Use `jira-onprem` for Server/Data Center (uses API v2 with wiki markup).
- **Description formatting** — For Server instances, always use Jira wiki markup for `--body` values. For Cloud, orbit auto-converts plain text to ADF. See `references/wiki-markup.md` for wiki markup syntax.
- **Field management** — `field create`, `field context-list`, `field option-list`, and `field option-add` are Cloud-only features.
- **Screen management** — Adding fields to screens makes them appear on issue create/edit forms. Use `tab-create` to group related fields into a dedicated tab.
- **Epic type cannot use `--parent` flag** — Jira rejects it because Epic is not a sub-task type. Use the `--field "customfield_27521=KEY"` (Parent Link) instead.
- **1Password integration** — Credentials in config can use `op://vault/item/field` and are resolved at runtime. Run `orbit auth` once to resolve and cache all secrets for 8 hours (single biometric prompt). Use `orbit auth clear` to wipe the cache. Without `orbit auth`, secrets are still resolved on each command but may trigger repeated biometric prompts.
