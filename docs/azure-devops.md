# Orbit Azure DevOps Command Reference

Manage Azure DevOps (VSTS) projects, work items, and saved queries from the command line.

**Top-level command:** `orbit azure-devops` (alias: `ado`)

**Persistent flag (all subcommands):**

| Flag | Description |
|------|-------------|
| `--service` | Azure DevOps service name, required only when a profile has multiple azure-devops services configured |

**Notes:**

- Supports both modern (`dev.azure.com/{org}`) and legacy (`{org}.visualstudio.com`) URLs.
- Authentication: basic auth with email as username and PAT as password.
- All commands support `-o json` and `-o yaml` for structured output.
- **Proxy bypass:** set `proxy: none` on the service to skip a profile-level proxy.
- **API version:** 7.1 (latest stable) appended to all requests.

---

## Table of Contents

- [version](#version)
- [project list](#project-list)
- [work-item (wi)](#work-item)
  - [wi view](#wi-view)
  - [wi list](#wi-list)
  - [wi create](#wi-create)
  - [wi update](#wi-update)
- [query](#query)
  - [query run](#query-run)
  - [query list](#query-list)

---

## version

Verify connectivity and show the authenticated user.

```
orbit ado version -p myprofile
```

---

## project list

List all projects in the Azure DevOps organization.

```
orbit ado project list -p myprofile
```

---

## work-item

Manage work items (bugs, tasks, user stories, epics, features).

**Alias:** `wi`

### wi view

View a single work item with all its fields.

| Argument | Position | Description |
|----------|----------|-------------|
| `id` | 1 | Work item ID (integer). |

```bash
orbit ado wi view 12345 -p myprofile
orbit ado wi view 12345 -o json -p myprofile
```

### wi list

Query work items using WIQL or convenience filter flags.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | config `default_project` | Project name. |
| `--wiql` | string | | Raw WIQL query (overrides filter flags). |
| `--type` | string | | Filter by work item type (Bug, Task, User Story, Epic, Feature). |
| `--state` | string | | Filter by state (New, Active, Resolved, Closed). |
| `--assigned-to` | string | | Filter by assigned user (email or display name). |
| `--max-results` | int | 50 | Maximum work items to fetch. |

```bash
# Filter by type and state
orbit ado wi list --project Fusion --type Bug --state New -p myprofile

# Raw WIQL
orbit ado wi list --project Fusion --wiql "SELECT [System.Id], [System.Title] FROM WorkItems WHERE [System.State] = 'Active'" -p myprofile

# Assigned to me
orbit ado wi list --project Fusion --assigned-to "me@company.com" -p myprofile
```

### wi create

Create a new work item using `--field` flags.

| Flag | Type | Description |
|------|------|-------------|
| `--project` | string | Project name. |
| `--type` | string | Work item type (required): Bug, Task, User Story, Epic, Feature. |
| `--field` | string (repeatable) | Field to set: `System.Title=My title`. |

```bash
orbit ado wi create --project Fusion --type Bug \
  --field "System.Title=Login page broken" \
  --field "Microsoft.VSTS.Common.Priority=1" -p myprofile
```

### wi update

Update fields on an existing work item.

| Argument | Position | Description |
|----------|----------|-------------|
| `id` | 1 | Work item ID. |

| Flag | Type | Description |
|------|------|-------------|
| `--field` | string (repeatable) | Field to set: `System.State=Resolved`. |

```bash
orbit ado wi update 12345 --field "System.State=Active" -p myprofile
orbit ado wi update 12345 \
  --field "System.State=Resolved" \
  --field "Microsoft.VSTS.Common.ResolvedReason=Fixed" -p myprofile
```

---

## query

Manage and execute saved queries.

### query run

Execute a saved query by its UUID and display the work items.

| Argument | Position | Description |
|----------|----------|-------------|
| `query-id` | 1 | The UUID of the saved query. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | config `default_project` | Project name. |
| `--max-results` | int | 50 | Maximum work items to fetch from results. |

```bash
orbit ado query run a670954b-e739-47ca-a09a-acc10f623123 --project Fusion -p myprofile
orbit ado query run a670954b-e739-47ca-a09a-acc10f623123 --project Fusion -o json -p myprofile
```

### query list

List saved queries and folders in a project.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | string | config `default_project` | Project name. |
| `--depth` | int | 1 | Folder tree depth. |

```bash
orbit ado query list --project Fusion -p myprofile
orbit ado query list --project Fusion --depth 3 -p myprofile
```

---

## Configuration Example

```yaml
profiles:
  epsilon:
    services:
      - name: azure-devops
        type: azure-devops
        base_url: https://vstsepsilon.visualstudio.com
        proxy: none                               # bypass profile-level proxy
        auth:
          method: basic
          username: user@company.com
          password: "YOUR_PAT_HERE"               # or op://vault/item/field
        options:
          default_project: Fusion                 # optional, saves typing --project
```

---

## WIQL Reference

Work Item Query Language (WIQL) is SQL-like:

```sql
SELECT [System.Id], [System.Title], [System.State]
FROM WorkItems
WHERE [System.TeamProject] = 'Fusion'
  AND [System.WorkItemType] = 'Bug'
  AND [System.State] <> 'Closed'
  AND [System.AssignedTo] = @Me
ORDER BY [System.CreatedDate] DESC
```

| WIQL | JQL Equivalent | Notes |
|------|---------------|-------|
| `[System.State] = 'Active'` | `status = "Active"` | Brackets + single quotes |
| `[System.WorkItemType] = 'Bug'` | `issuetype = Bug` | |
| `@Me` | `currentUser()` | Current authenticated user |
| `<>` | `!=` | Not-equal operator |
| `IN ('A', 'B')` | `in (A, B)` | Same concept |
