# logx Quick Reference

## Installation

```bash
# Clone and build
git clone https://github.com/jatsandaruwan/logx.git
cd logx
make build

# Or use install script
chmod +x install.sh
./install.sh
```

## Basic Commands

### User Management
```bash
logx user add                    # Add new user
logx user list                   # List all users
logx user delete <name>          # Delete user
```

### App Management
```bash
logx app add                     # Add new app
logx app list                    # List all apps
logx app update <name>           # Update app
logx app delete <name>           # Delete app
```

### View Logs
```bash
logx <appname>                   # View current logs
logx <appname> 2025-09-10        # View logs for date
logx <appname> --server <IP>     # View from specific server
```

### Editor Configuration
```bash
logx editor set <command>        # Set custom editor
logx editor show                 # Show current editor
```

### Help
```bash
logx help                        # Show help
logx version                     # Show version
```

## Configuration Examples

### Add User
```
User Name: prod-admin
SSH Username: root
SSH Password: ********
```

### Add App
```
App Name: testapp
User: prod-admin
Log Path: /logs/testapp/testapp.log
Pattern: testapp-{date}.log
Date Format: 2006-01-02
Servers: 192.168.0.1, 192.168.0.2
```

## Date Format Patterns

| Pattern | Example Output | Use Case |
|---------|----------------|----------|
| `2006-01-02` | 2025-09-10 | ISO format |
| `20060102` | 20250910 | Compact |
| `02-01-2006` | 10-09-2025 | DD-MM-YYYY |
| `01/02/2006` | 09/10/2025 | US format |
| `2006_01_02` | 2025_09_10 | Underscore |

## Log Patterns

| Pattern | Date Format | Result |
|---------|-------------|--------|
| `app-{date}.log` | `2006-01-02` | `app-2025-09-10.log` |
| `app.log-{date}` | `20060102` | `app.log-20250910` |
| `{date}_app.log` | `2006-01-02` | `2025-09-10_app.log` |
| `app_{date}.log` | `2006_01_02` | `app_2025_09_10.log` |

## File Locations

### Config
- Linux/macOS: `~/.config/logx/config.xml`
- Windows: `%APPDATA%\logx\config.xml`

### Credentials
- Windows: Windows Credential Manager
- macOS: Keychain
- Linux: Secret Service (gnome-keyring)

## Common Workflows

### Setup New Application
```bash
# 1. Add user (if not exists)
logx user add

# 2. Add application
logx app add

# 3. Test connection
logx myapp
```

### Daily Log Viewing
```bash
# Today's logs
logx myapp

# Yesterday's logs
logx myapp 2025-09-29

# Specific date
logx myapp 2025-09-10
```

### Troubleshooting
```bash
# Check configuration
logx app list
logx user list

# Test specific server
logx myapp --server 192.168.0.1

# Check editor setting
logx editor show
```

## Editor Commands

### Common Editors
```bash
# VS Code
logx editor set code

# Notepad++
logx editor set notepad++

# Sublime Text
logx editor set subl

# Vim
logx editor set vim

# Nano
logx editor set nano

# VS Code with wait flag
logx editor set "code -w"
```

## Error Messages

| Error | Solution |
|-------|----------|
| "User not found" | Run `logx user add` |
| "App not found" | Run `logx app add` |
| "Failed to connect" | Check IP, credentials, SSH service |
| "Log file not found" | Verify path and date format |
| "No suitable editor" | Run `logx editor set <editor>` |

## Tips & Tricks

### 1. Multiple Servers
Add all servers for an app to download from all at once:
```bash
# During app add, enter multiple IPs
Server IP: 192.168.0.1
Server IP: 192.168.0.2
Server IP: 192.168.0.3
Server IP: [press Enter to finish]
```

### 2. Quick Access
Create shell aliases:
```bash
# In ~/.bashrc or ~/.zshrc
alias logs-prod='logx prodapp'
alias logs-dev='logx devapp'
```

### 3. Date Shortcuts
```bash
# Use date command for calculations
logx myapp $(date -d "yesterday" +%Y-%m-%d)
logx myapp $(date -d "7 days ago" +%Y-%m-%d)
```

### 4. Batch Operations
```bash
# View logs from multiple apps
for app in app1 app2 app3; do
  logx $app 2025-09-10
done
```

### 5. Editor Integration
```bash
# Open multiple logs and compare
logx app1 2025-09-10 &
logx app2 2025-09-10 &
wait
```

## Security Checklist

- ✅ Use strong SSH passwords
- ✅ Regularly update credentials
- ✅ Limit server access to necessary IPs
- ✅ Review user/app list periodically
- ✅ Protect config file permissions (0600)
- ✅ Don't share config file (credentials in keyring)
- ✅ Use SSH keys when possible (future feature)

## Build from Source

```bash
# Get dependencies
go mod download

# Build for current platform
go build -o logx cmd/logx/main.go

# Build for all platforms
make build-all

# Install
sudo make install
```

## Keyboard Shortcuts

During interactive prompts:
- `Ctrl+C` or `Esc` - Cancel/Exit
- `Enter` - Confirm/Next
- `Backspace` - Delete character

## Configuration XML Structure

```xml
<config>
  <users>
    <user id="ID" name="Name" username="SSHUser"/>
  </users>
  <apps>
    <app name="AppName">
      <user-ref>UserID</user-ref>
      <log-path>/path/to/log</log-path>
      <log-pattern>pattern-{date}.log</log-pattern>
      <date-format>2006-01-02</date-format>
      <servers>
        <server>IP1</server>
        <server>IP2</server>
      </servers>
    </app>
  </apps>
  <editor>command</editor>
</config>
```

## Quick Diagnostics

```bash
# Check version
logx version

# List configuration
logx user list
logx app list

# Test SSH manually
ssh username@192.168.0.1

# Check log file manually
ssh username@192.168.0.1 'ls -l /logs/app/app-2025-09-10.log'

# Verify editor
which code
which notepad++
```

## Support

- Documentation: [README.md](README.md)
- Setup Guide: [SETUP.md](SETUP.md)
- Issues: https://github.com/jat-sandaruwan/logx/issues