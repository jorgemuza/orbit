# QMetry Command Reference

Complete reference for all `orbit qmetry` (alias: `qm`) commands.

## testcase (alias: tc)

### `qmetry testcase create`

Create a test case.

| Flag | Required | Description |
|------|----------|-------------|
| `--summary` | Yes | Test case summary |
| `--project-id` | Yes | QMetry project ID |
| `--description` | No | Description text |
| `--precondition` | No | Precondition text |
| `--folder-id` | No | Target folder ID |
| `--steps` | No | BDD steps as JSON array |

Steps JSON format:
```json
[
  {"stepDetails": "Given ...", "testData": "", "expectedResult": ""},
  {"stepDetails": "When ...", "testData": "input data", "expectedResult": ""},
  {"stepDetails": "Then ...", "testData": "", "expectedResult": "expected outcome"}
]
```

### `qmetry testcase view <key>`

View test case details including steps.

### `qmetry testcase list`

List test cases in a project.

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--project-id` | Yes | | QMetry project ID |
| `--folder-id` | No | | Filter by folder |
| `--limit` | No | 25 | Max results |

---

## testcycle (alias: cycle)

### `qmetry testcycle create`

Create a test cycle.

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Cycle name |
| `--project-id` | Yes | QMetry project ID |
| `--description` | No | Cycle description |
| `--folder-id` | No | Target folder ID |
| `--start-date` | No | Start date (YYYY-MM-DD) |
| `--end-date` | No | End date (YYYY-MM-DD) |

### `qmetry testcycle list`

List test cycles in a project.

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--project-id` | Yes | | QMetry project ID |
| `--limit` | No | 25 | Max results |

---

## project

List QMetry projects.

```
orbit qm project
```

---

## folder

List folders in a QMetry project.

| Flag | Required | Description |
|------|----------|-------------|
| `--project-id` | Yes | QMetry project ID |
| `--type` | No | Folder type: `testcase`, `testcycle` |
