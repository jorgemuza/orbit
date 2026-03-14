# orbit

AI Development Lifecycle CLI — manage connections to Jira, Confluence, GitLab, Bitbucket, and GitHub through named profiles.

## Install

### Homebrew (macOS / Linux)

```bash
brew install jorgemuza/tap/orbit
```

### Scoop (Windows)

```powershell
scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket
scoop install orbit
```

### Manual

Download the latest release from the [Releases](https://github.com/jorgemuza/orbit/releases) page.

## Usage

```bash
# Configure a profile
orbit profile create myprofile

# Test connections
orbit service ping -p myprofile

# Jira, Confluence, GitLab, Bitbucket, GitHub commands
orbit jira issue view PROJ-123
orbit confluence page view --title "My Page"
orbit gitlab variable list --project my/project
orbit bitbucket pr list
orbit github pr list
```

## Documentation

See the [docs](https://github.com/jorgemuza/orbit/wiki) for full usage details.

## License

MIT
