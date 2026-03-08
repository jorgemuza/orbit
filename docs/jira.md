# Orbit Jira Command Reference

Complete reference for all `orbit jira` CLI commands.

## Table of Contents

- [Global Flags](#global-flags)
- [Issue Commands](#issue-commands)
  - [issue list](#issue-list)
  - [issue view](#issue-view)
  - [issue create](#issue-create)
  - [issue edit](#issue-edit)
  - [issue assign](#issue-assign)
  - [issue move](#issue-move)
  - [issue delete](#issue-delete)
  - [issue comment](#issue-comment)
  - [issue link](#issue-link)
  - [issue unlink](#issue-unlink)
  - [issue worklog](#issue-worklog)
  - [issue clone](#issue-clone)
- [Epic Commands](#epic-commands)
  - [epic list](#epic-list)
  - [epic create](#epic-create)
  - [epic add](#epic-add)
  - [epic remove](#epic-remove)
- [Sprint Commands](#sprint-commands)
  - [sprint list](#sprint-list)
  - [sprint add](#sprint-add)
- [Board Commands](#board-commands)
  - [board list](#board-list)
- [Project Commands](#project-commands)
  - [project list](#project-list)
- [Release Commands](#release-commands)
  - [release list](#release-list)
- [Field Commands](#field-commands)
  - [field list](#field-list)
  - [field create](#field-create)
  - [field context-list](#field-context-list)
  - [field option-list](#field-option-list)
  - [field option-add](#field-option-add)
- [Screen Commands](#screen-commands)
  - [screen list](#screen-list)
  - [screen tab-list](#screen-tab-list)
  - [screen tab-create](#screen-tab-create)
  - [screen field-list](#screen-field-list)
  - [screen field-add](#screen-field-add)
  - [screen field-remove](#screen-field-remove)
  - [screen field-move](#screen-field-move)
- [Status Commands](#status-commands)
  - [status list](#status-list)
- [Issue Type Commands](#issue-type-commands)
  - [issuetype-list](#issuetype-list)
- [Notes](#notes)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--service` | Jira service name. Required only when a profile has multiple Jira services configured. |

All examples below use `-p myprofile` to specify the orbit profile.

---

## Issue Commands

### issue list

List and search issues. Aliases: `ls`, `search`.

```
orbit jira issue list [flags] -p myprofile
```

#### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--jql` | `-q` | | Raw JQL query string |
| `--project` | | | Filter by project key |
| `--type` | `-t` | | Filter by issue type (e.g. Bug, Story, Task) |
| `--status` | `-s` | | Filter by status (repeatable) |
| `--priority` | `-y` | | Filter by priority |
| `--assignee` | `-a` | | Filter by assignee username |
| `--reporter` | `-r` | | Filter by reporter username |
| `--label` | `-l` | | Filter by label (repeatable) |
| `--component` | `-C` | | Filter by component |
| `--parent` | `-P` | | Filter by parent issue key |
| `--created` | | | Filter by created date |
| `--updated` | | | Filter by updated date |
| `--created-after` | | | Issues created after this date |
| `--created-before` | | | Issues created before this date |
| `--updated-after` | | | Issues updated after this date |
| `--updated-before` | | | Issues updated before this date |
| `--order-by` | | `created` | Field to order results by |
| `--reverse` | | | Reverse sort order |
| `--start-at` | | | Pagination offset |
| `--max-results` | | `50` | Maximum number of results to return |

#### Examples

```bash
# List all issues in a project
orbit jira issue list --project MYPROJ -p myprofile

# Search with raw JQL
orbit jira issue list -q "project = MYPROJ AND status = 'In Progress'" -p myprofile

# Filter by multiple statuses
orbit jira issue list --project MYPROJ -s "To Do" -s "In Progress" -p myprofile

# Filter by assignee and type
orbit jira issue list --project MYPROJ -a john.doe -t Bug -p myprofile

# Issues created after a date, ordered by priority
orbit jira issue list --project MYPROJ --created-after 2026-01-01 --order-by priority -p myprofile

# Filter by label and component
orbit jira issue list --project MYPROJ -l backend -l urgent -C "API Team" -p myprofile

# Reverse sort with limited results
orbit jira issue list --project MYPROJ --order-by updated --reverse --max-results 10 -p myprofile
```

---

### issue view

View a single issue's details. Aliases: `show`.

```
orbit jira issue view [issue-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key (e.g. MYPROJ-123) |

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--comments` | `1` | Number of comments to display |

#### Examples

```bash
# View an issue with the default 1 comment
orbit jira issue view MYPROJ-123 -p myprofile

# View an issue with the last 5 comments
orbit jira issue view MYPROJ-123 --comments 5 -p myprofile
```

---

### issue create

Create a new issue.

```
orbit jira issue create [flags] -p myprofile
```

#### Flags

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--project` | | Yes | Project key |
| `--type` | `-t` | Yes | Issue type (e.g. Bug, Story, Task, Epic) |
| `--summary` | `-s` | Yes | Issue summary/title |
| `--body` | `-b` | | Issue description |
| `--priority` | `-y` | | Priority level (e.g. High, Medium, Low) |
| `--assignee` | `-a` | | Assignee username |
| `--reporter` | `-r` | | Reporter username |
| `--parent` | `-P` | | Parent issue key (for subtasks) |
| `--label` | `-l` | | Label (repeatable) |
| `--component` | `-C` | | Component name (repeatable) |
| `--fix-version` | | | Fix version (repeatable) |
| `--affects-version` | | | Affects version (repeatable) |
| `--original-estimate` | `-e` | | Original time estimate (e.g. 2h, 1d, 1w) |
| `--epic-name` | | | Epic name (when creating an Epic) |
| `--field` | `-F` | | Set arbitrary field as key=value (repeatable) |

#### Examples

```bash
# Create a basic task
orbit jira issue create --project MYPROJ -t Task -s "Implement login page" -p myprofile

# Create a bug with description and priority
orbit jira issue create --project MYPROJ -t Bug \
  -s "Login button unresponsive on mobile" \
  -b "The login button does not respond to taps on iOS Safari." \
  -y High -p myprofile

# Create a story with labels, components, and an estimate
orbit jira issue create --project MYPROJ -t Story \
  -s "User password reset flow" \
  -l backend -l security \
  -C "Auth Service" \
  -e 3d -p myprofile

# Create a subtask under a parent issue
orbit jira issue create --project MYPROJ -t Sub-task \
  -s "Write unit tests for login" \
  -P MYPROJ-100 -p myprofile

# Create an issue with custom fields
orbit jira issue create --project MYPROJ -t Story \
  -s "Add metrics dashboard" \
  -F "customfield_10010=team-alpha" \
  -F "customfield_10020=Q1-2026" -p myprofile

# Create an issue with fix and affects versions
orbit jira issue create --project MYPROJ -t Bug \
  -s "Crash on startup" \
  --fix-version "2.1.0" \
  --affects-version "2.0.0" --affects-version "2.0.1" -p myprofile
```

---

### issue edit

Edit an existing issue. Aliases: `update`.

```
orbit jira issue edit [issue-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key to edit |

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--summary` | `-s` | New summary/title |
| `--body` | `-b` | New description |
| `--priority` | `-y` | New priority |
| `--label` | `-l` | Set labels (repeatable, replaces existing) |
| `--component` | `-C` | Set components (repeatable, replaces existing) |
| `--fix-version` | | Set fix versions (repeatable, replaces existing) |
| `--field` | `-F` | Set arbitrary field as key=value (repeatable) |

#### Examples

```bash
# Update summary
orbit jira issue edit MYPROJ-123 -s "Updated title" -p myprofile

# Update description and priority
orbit jira issue edit MYPROJ-123 -b "New description text" -y Critical -p myprofile

# Set labels and a custom field
orbit jira issue edit MYPROJ-123 -l urgent -l production -F "customfield_10010=new-value" -p myprofile
```

---

### issue assign

Assign an issue to a user.

```
orbit jira issue assign [issue-key] [assignee] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key |
| `assignee` | Username to assign. Use `x` to unassign. |

#### Examples

```bash
# Assign to a user
orbit jira issue assign MYPROJ-123 john.doe -p myprofile

# Unassign
orbit jira issue assign MYPROJ-123 x -p myprofile
```

---

### issue move

Transition an issue to a new status. Aliases: `transition`.

```
orbit jira issue move [issue-key] [state] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key |
| `state` | Target status/transition name (e.g. "In Progress", "Done") |

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--comment` | | Add a comment during the transition |
| `--resolution` | `-R` | Set resolution (e.g. Fixed, Won't Fix) |

#### Examples

```bash
# Move to In Progress
orbit jira issue move MYPROJ-123 "In Progress" -p myprofile

# Move to Done with a resolution and comment
orbit jira issue move MYPROJ-123 Done -R Fixed --comment "Verified in staging" -p myprofile
```

---

### issue delete

Delete an issue. Aliases: `rm`, `remove`.

```
orbit jira issue delete [issue-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key to delete |

#### Flags

| Flag | Description |
|------|-------------|
| `--cascade` | Also delete all subtasks |

#### Examples

```bash
# Delete a single issue
orbit jira issue delete MYPROJ-123 -p myprofile

# Delete an issue and all its subtasks
orbit jira issue delete MYPROJ-100 --cascade -p myprofile
```

---

### issue comment

Add a comment to an issue.

```
orbit jira issue comment [issue-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key |

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--body` | `-b` | Comment body text |

#### Examples

```bash
orbit jira issue comment MYPROJ-123 -b "This has been deployed to staging." -p myprofile
```

---

### issue link

Create a link between two issues.

```
orbit jira issue link [inward-key] [outward-key] [link-type] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `inward-key` | The inward issue key (e.g. MYPROJ-100) |
| `outward-key` | The outward issue key (e.g. MYPROJ-200) |
| `link-type` | Link type name (e.g. "Blocks", "Relates", "Cloners") |

#### Examples

```bash
# MYPROJ-100 blocks MYPROJ-200
orbit jira issue link MYPROJ-100 MYPROJ-200 Blocks -p myprofile
```

---

### issue unlink

Remove a link between two issues.

```
orbit jira issue unlink [inward-key] [outward-key] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `inward-key` | The inward issue key |
| `outward-key` | The outward issue key |

#### Examples

```bash
orbit jira issue unlink MYPROJ-100 MYPROJ-200 -p myprofile
```

---

### issue worklog

Log time spent on an issue.

```
orbit jira issue worklog [issue-key] [time-spent] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key |
| `time-spent` | Time spent (e.g. 2h, 1d, 30m) |

#### Flags

| Flag | Description |
|------|-------------|
| `--comment` | Comment describing the work done |

#### Examples

```bash
# Log 2 hours
orbit jira issue worklog MYPROJ-123 2h -p myprofile

# Log 1 day with a comment
orbit jira issue worklog MYPROJ-123 1d --comment "Completed code review" -p myprofile
```

---

### issue clone

Clone an existing issue.

```
orbit jira issue clone [issue-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The issue key to clone |

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--summary` | `-s` | Override the cloned issue's summary |
| `--replace` | `-H` | Find-and-replace in cloned content as `find:replace` (repeatable) |

#### Examples

```bash
# Clone an issue
orbit jira issue clone MYPROJ-123 -p myprofile

# Clone with a new summary
orbit jira issue clone MYPROJ-123 -s "Cloned: new feature variant" -p myprofile

# Clone with text replacements
orbit jira issue clone MYPROJ-123 -H "v1:v2" -H "staging:production" -p myprofile
```

---

## Epic Commands

### epic list

List epics or list issues within an epic. Aliases: `ls`.

```
orbit jira epic list [epic-key] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `epic-key` | (Optional) Epic key to list its child issues. Omit to list all epics. |

#### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | | | Filter by project key |
| `--status` | `-s` | | Filter by status (repeatable) |
| `--max-results` | | `50` | Maximum number of results |

#### Examples

```bash
# List all epics in a project
orbit jira epic list --project MYPROJ -p myprofile

# List issues inside a specific epic
orbit jira epic list MYPROJ-50 -p myprofile

# List open issues in an epic
orbit jira epic list MYPROJ-50 -s "To Do" -s "In Progress" -p myprofile
```

---

### epic create

Create a new epic.

```
orbit jira epic create [flags] -p myprofile
```

#### Flags

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--project` | | Yes | Project key |
| `--name` | `-n` | Yes | Epic name |
| `--summary` | `-s` | | Epic summary/title |
| `--body` | `-b` | | Epic description |
| `--priority` | `-y` | | Priority level |
| `--label` | `-l` | | Label (repeatable) |
| `--component` | `-C` | | Component name (repeatable) |

#### Examples

```bash
# Create an epic
orbit jira epic create --project MYPROJ -n "Q1 Auth Overhaul" \
  -s "Revamp authentication system" \
  -b "Migrate to OAuth2 and add MFA support." \
  -y High -p myprofile

# Create an epic with labels and components
orbit jira epic create --project MYPROJ -n "Performance" \
  -l performance -l infrastructure \
  -C "Platform" -p myprofile
```

---

### epic add

Add issues to an epic.

```
orbit jira epic add [epic-key] [issue-keys...] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `epic-key` | The epic key to add issues to |
| `issue-keys` | One or more issue keys to add (max 50) |

#### Examples

```bash
orbit jira epic add MYPROJ-50 MYPROJ-101 MYPROJ-102 MYPROJ-103 -p myprofile
```

---

### epic remove

Remove issues from their epic.

```
orbit jira epic remove [issue-keys...] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `issue-keys` | One or more issue keys to remove from their epic (max 50) |

#### Examples

```bash
orbit jira epic remove MYPROJ-101 MYPROJ-102 -p myprofile
```

---

## Sprint Commands

### sprint list

List sprints or list issues within a sprint. Aliases: `ls`.

```
orbit jira sprint list [sprint-id] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `sprint-id` | (Optional) Sprint ID to list its issues. Omit to list sprints. |

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--board-id` | | Board ID to list sprints for |
| `--state` | | Filter by sprint state: `future`, `active`, or `closed` |
| `--max-results` | `50` | Maximum number of results |

#### Examples

```bash
# List sprints for a board
orbit jira sprint list --board-id 42 -p myprofile

# List only active sprints
orbit jira sprint list --board-id 42 --state active -p myprofile

# List issues in a sprint
orbit jira sprint list 315 -p myprofile
```

---

### sprint add

Add issues to a sprint.

```
orbit jira sprint add [sprint-id] [issue-keys...] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `sprint-id` | The target sprint ID |
| `issue-keys` | One or more issue keys to add (max 50) |

#### Examples

```bash
orbit jira sprint add 315 MYPROJ-101 MYPROJ-102 MYPROJ-103 -p myprofile
```

---

## Board Commands

### board list

List boards. Aliases: `ls`.

```
orbit jira board list [flags] -p myprofile
```

#### Flags

| Flag | Description |
|------|-------------|
| `--project` | Filter boards by project key |

#### Examples

```bash
# List all boards
orbit jira board list -p myprofile

# List boards for a specific project
orbit jira board list --project MYPROJ -p myprofile
```

---

## Project Commands

### project list

List all accessible projects. Aliases: `ls`.

```
orbit jira project list -p myprofile
```

#### Examples

```bash
orbit jira project list -p myprofile
```

---

## Release Commands

### release list

List versions/releases for a project. Aliases: `ls`.

```
orbit jira release list [flags] -p myprofile
```

#### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--project` | Yes | Project key |

#### Examples

```bash
orbit jira release list --project MYPROJ -p myprofile
```

---

## Field Commands

### field list

List available fields.

```
orbit jira field list [flags] -p myprofile
```

#### Flags

| Flag | Description |
|------|-------------|
| `--filter` | Filter fields by name substring |
| `--custom` | Show only custom fields (boolean) |

#### Examples

```bash
# List all fields
orbit jira field list -p myprofile

# List only custom fields
orbit jira field list --custom -p myprofile

# Search for fields by name
orbit jira field list --filter "sprint" -p myprofile
```

---

### field create

Create a custom field. **Cloud only.**

```
orbit jira field create [flags] -p myprofile
```

#### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Field name |
| `--type` | Yes | Field type. Shorthands: `select`, `multiselect`, `number`/`float`, `checkbox`/`checkboxes`, `text`, `textarea` |
| `--description` | | Field description |
| `--searcher` | | Searcher key for the field |

#### Examples

```bash
# Create a select field
orbit jira field create --name "Environment" --type select --description "Deployment environment" -p myprofile

# Create a numeric field
orbit jira field create --name "Story Points" --type number -p myprofile

# Create a multi-select field
orbit jira field create --name "Affected Teams" --type multiselect -p myprofile
```

---

### field context-list

List contexts for a custom field. **Cloud only.**

```
orbit jira field context-list [field-id] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `field-id` | The custom field ID (e.g. customfield_10100) |

#### Examples

```bash
orbit jira field context-list customfield_10100 -p myprofile
```

---

### field option-list

List options for a custom field context. **Cloud only.**

```
orbit jira field option-list [field-id] [context-id] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `field-id` | The custom field ID |
| `context-id` | The context ID |

#### Examples

```bash
orbit jira field option-list customfield_10100 10200 -p myprofile
```

---

### field option-add

Add options to a custom field context. **Cloud only.**

```
orbit jira field option-add [field-id] [context-id] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `field-id` | The custom field ID |
| `context-id` | The context ID |

#### Flags

| Flag | Description |
|------|-------------|
| `--values` | Comma-separated list of option values to add |

#### Examples

```bash
orbit jira field option-add customfield_10100 10200 --values "staging,production,development" -p myprofile
```

---

## Screen Commands

### screen list

List screens.

```
orbit jira screen list [flags] -p myprofile
```

#### Flags

| Flag | Description |
|------|-------------|
| `--filter` | Filter screens by name substring |

#### Examples

```bash
# List all screens
orbit jira screen list -p myprofile

# Filter by name
orbit jira screen list --filter "Default" -p myprofile
```

---

### screen tab-list

List tabs on a screen.

```
orbit jira screen tab-list [screen-id] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |

#### Examples

```bash
orbit jira screen tab-list 1 -p myprofile
```

---

### screen tab-create

Create a new tab on a screen.

```
orbit jira screen tab-create [screen-id] [tab-name] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |
| `tab-name` | Name for the new tab |

#### Examples

```bash
orbit jira screen tab-create 1 "Custom Fields" -p myprofile
```

---

### screen field-list

List fields on a screen tab.

```
orbit jira screen field-list [screen-id] [tab-id] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |
| `tab-id` | The tab ID |

#### Examples

```bash
orbit jira screen field-list 1 10001 -p myprofile
```

---

### screen field-add

Add fields to a screen tab.

```
orbit jira screen field-add [screen-id] [tab-id] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |
| `tab-id` | The tab ID |

#### Flags

| Flag | Description |
|------|-------------|
| `--fields` | Comma-separated list of field IDs to add |

#### Examples

```bash
orbit jira screen field-add 1 10001 --fields "customfield_10100,customfield_10101" -p myprofile
```

---

### screen field-remove

Remove fields from a screen tab.

```
orbit jira screen field-remove [screen-id] [tab-id] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |
| `tab-id` | The tab ID |

#### Flags

| Flag | Description |
|------|-------------|
| `--fields` | Comma-separated list of field IDs to remove |

#### Examples

```bash
orbit jira screen field-remove 1 10001 --fields "customfield_10100" -p myprofile
```

---

### screen field-move

Move fields between screen tabs.

```
orbit jira screen field-move [screen-id] [source-tab-id] [target-tab-id] [flags] -p myprofile
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `screen-id` | The screen ID |
| `source-tab-id` | Tab ID to move fields from |
| `target-tab-id` | Tab ID to move fields to |

#### Flags

| Flag | Description |
|------|-------------|
| `--fields` | Comma-separated list of field IDs to move |

#### Examples

```bash
orbit jira screen field-move 1 10001 10002 --fields "customfield_10100,customfield_10101" -p myprofile
```

---

## Status Commands

### status list

List all available issue statuses.

```
orbit jira status list -p myprofile
```

#### Examples

```bash
orbit jira status list -p myprofile
```

---

## Issue Type Commands

### issuetype-list

List all available issue types. This is a top-level subcommand under `jira`.

```
orbit jira issuetype-list -p myprofile
```

#### Examples

```bash
orbit jira issuetype-list -p myprofile
```

---

## Notes

### Cloud vs Server

Orbit supports both Jira Cloud (API v3) and Jira Server/Data Center (API v2). Key differences:

| Feature | Cloud (API v3) | Server (API v2) |
|---------|----------------|-----------------|
| Descriptions | ADF (Atlassian Document Format) | Wiki markup |
| `field create` | Supported | Not supported |
| `field context-list` | Supported | Not supported |
| `field option-list` | Supported | Not supported |
| `field option-add` | Supported | Not supported |
| Screen commands | Supported | Supported |

When using Cloud, orbit automatically converts plain text descriptions to ADF format.

### Custom Fields

The `--field` / `-F` flag on `issue create` and `issue edit` allows setting arbitrary fields by their field ID:

```bash
orbit jira issue create --project MYPROJ -t Task -s "My task" \
  -F "customfield_10010=my-value" \
  -F "customfield_10020=another-value" -p myprofile
```

Use `orbit jira field list` to discover field IDs.

### Multiple Jira Services

If your orbit profile is configured with multiple Jira services, use the `--service` flag to specify which one:

```bash
orbit jira issue list --project MYPROJ --service my-jira-cloud -p myprofile
```
