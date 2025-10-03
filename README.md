# logx - Remote Log Viewer CLI

A command-line tool for viewing and managing log files from remote servers via SSH. Built with Go and Bubble Tea.

## Features

- 🔐 **Secure Credential Storage**: Uses system keyring (Windows Credential Manager, macOS Keychain, Linux Secret Service)
- 🖥️ **Multi-Server Support**: Access logs from multiple servers with the same configuration
- 📅 **Date-Based Log Rolling**: Support for various date formats in log filenames
- ✏️ **Configurable Editors**: Use your preferred text editor (Notepad++, VS Code, vim, etc.)
- 🎨 **Interactive TUI**: Beautiful terminal UI using Bubble Tea
- 📝 **XML Configuration**: Easy-to-read configuration format

## Installation

### Prerequisites

- Go 1.21 or higher
- SSH access to remote servers

### Build from Source

```bash
git clone https://github.com/jatsandaruwan/logx.git
cd logx
go build -o logx cmd/logx/main.go
```

### Install Globally

```bash
go install github.com/jatsandaruwan/logx/cmd/logx@latest
```

Or copy the binary to your PATH:

```bash
# Linux/macOS
sudo cp logx /usr/local/bin/

# Windows
# Copy logx.exe to C:\Windows\System32\ or add to PATH
```

## Quick Start

### 1. Add a User

First, add SSH credentials:

```bash
logx user add
```

You'll be prompted for:
- User name (identifier)
- SSH username
- SSH password

### 2. Add an Application

Configure an application's log location:

```bash
logx app add
```

You'll be prompted for:
- App name
- User (select from configured users)
- Log file path (e.g., `/logs/testapp/testapp.log`)
- Log filename pattern with `{date}` placeholder
  - Example: `testapp-{date}.log`
  - Example: `testapp.log-{date}`
- Date format (Go format)
  - `2006-01-02` → 2025-09-10
  - `20060102` → 20250910
  - `02-01-2006` → 10-09-2025
- Server IPs (one per line, empty to finish)

### 3. View Logs

View current logs:
```bash
logx testapp
```

View logs for a specific date:
```bash
logx testapp 2025-09-10
```

View logs from a specific server:
```bash
logx testapp --server 192.168.0.1
logx testapp 2025-09-10 --server 192.168.0.1
```

## Usage

### User Management

```bash
# Add a new user
logx user add

# List all users
logx user list

# Delete a user
logx user delete <username>
```

### Application Management

```bash
# Add a new application
logx app add

# List all applications
logx app list

# Update an application
logx app update <appname>

# Delete an application
logx app delete <appname>
```

### Editor Configuration

```bash
# Set custom editor
logx editor set "code"
logx editor set "notepad++"
logx editor set "vim"

# Show current editor
logx editor show
```

### Viewing Logs

```bash
# View current logs (all servers)
logx <appname>

# View logs for a specific date (all servers)
logx <appname> <YYYY-MM-DD>

# View logs from a specific server
logx <appname> --server <IP>
logx <appname> <YYYY-MM-DD> --server <IP>
```

## Configuration

### Configuration File Location

- **Linux/macOS**: `~/.config/logx/config.xml`
- **Windows**: `%APPDATA%\logx\config.xml`

### Example Configuration

```xml
<?xml version="1.0" encoding="UTF-8"?>
<config>
  <users>
    <user id="admin" name="admin" username="root"/>
    <user id="devuser" name="devuser" username="developer"/>
  </users>
  <apps>
    <app name="testapp">
      <user-ref>admin</user-ref>
      <log-path>/logs/testapp/testapp.log</log-path>
      <log-pattern>testapp-{date}.log</log-pattern>
      <date-format>2006-01-02</date-format>
      <servers>
        <server>192.168.0.1</server>
        <server>192.168.0.2</server>
      </servers>
    </app>
    <app name="webapp">
      <user-ref>devuser</user-ref>
      <log-path>/var/log/webapp/app.log</log-path>
      <log-pattern>app.log-{date}</log-pattern>
      <date-format>20060102</date-format>
      <servers>
        <server>10.0.0.5</server>
      </servers>
    </app>
  </apps>
  <editor>code</editor>
</config>
```

## Date Format Reference

Go uses a specific reference time for formatting: **Mon Jan 2 15:04:05 MST 2006**

Common patterns:

| Pattern | Output | Description |
|---------|--------|-------------|
| `2006-01-02` | 2025-09-10 | ISO 8601 format |
| `20060102` | 20250910 | Compact format |
| `02-01-2006` | 10-09-2025 | DD-MM-YYYY format |
| `01/02/2006` | 09/10/2025 | US format |
| `2006_01_02` | 2025_09_10 | Underscore format |

## Security

- **Credentials**: Stored securely in system keyring
  - Windows: Windows Credential Manager
  - macOS: Keychain
  - Linux: Secret Service API (GNOME Keyring, KWallet)
- **SSH**: Password authentication (SSH keys support can be added)
- **Config File**: Stored with 0600 permissions (user read/write only)

## Platform Support

- ✅ Windows
- ✅ macOS
- ✅ Linux

## Default Editors by Platform

- **Windows**: Notepad++ → Notepad
- **macOS**: VS Code → Sublime Text → TextEdit
- **Linux**: VS Code → gedit → kate → nano → vim

## Troubleshooting

### SSH Connection Failures

- Verify server IP and port (default: 22)
- Check username and password
- Ensure SSH service is running on remote server
- Check firewall rules

### Log Files Not Found

- Verify log path is correct
- Check date format matches actual log filename
- Ensure log pattern includes `{date}` placeholder
- Verify user has read permissions on remote server

### Editor Not Opening

- Set a custom editor: `logx editor set <editor>`
- Ensure the editor is in your PATH
- For Windows Notepad++, verify installation path

## Development

### Project Structure

```
logx/
├── cmd/
│   └── logx/           # Main CLI entry point
│       └── main.go
├── internal/
│   ├── config/         # XML configuration
│   │   └── config.go
│   ├── editor/         # Editor handling
│   │   └── editor.go
│   ├── ssh/            # SSH connections
│   │   └── ssh.go
│   ├── ui/             # Bubble Tea UI
│   │   ├── app.go
│   │   └── user.go
│   ├── vault/          # Credential storage
│   │   └── vault.go
│   └── viewer/         # Log viewing
│       └── viewer.go
├── go.mod
├── go.sum
└── README.md
```

### Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/zalando/go-keyring` - System keyring access
- `golang.org/x/crypto/ssh` - SSH client
- `golang.org/x/term` - Terminal utilities

### Building

```bash
# Build for current platform
go build -o logx cmd/logx/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o logx-linux cmd/logx/main.go
GOOS=windows GOARCH=amd64 go build -o logx.exe cmd/logx/main.go
GOOS=darwin GOARCH=amd64 go build -o logx-darwin cmd/logx/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Author

Thilina Sandaruwan

## Support

For issues and feature requests, please create an issue on GitHub.