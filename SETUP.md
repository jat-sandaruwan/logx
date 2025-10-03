# logx Setup Guide

This guide will walk you through setting up and using logx for the first time.

## Prerequisites

Before installing logx, ensure you have:

1. **Go 1.21+** installed
2. **SSH access** to your remote servers
3. **SSH credentials** (username and password)
4. **Text editor** of your choice installed (optional)

## Installation Steps

### Step 1: Clone and Build

```bash
# Clone the repository
git clone https://github.com/jatsandaruwan/logx.git
cd logx

# Download dependencies
go mod download

# Build the binary
make build

# Or build manually
go build -o logx cmd/logx/main.go
```

### Step 2: Install System-Wide (Optional)

**Linux/macOS:**
```bash
sudo make install
# Or manually
sudo cp build/logx /usr/local/bin/
```

**Windows:**
```powershell
# Copy logx.exe to a directory in your PATH
# For example: C:\Windows\System32\
```

### Step 3: Verify Installation

```bash
logx version
# Should output: logx version 1.0.0
```

## Initial Configuration

### Step 1: Add Your First User

Users store SSH credentials securely in your system's keyring.

```bash
logx user add
```

You'll be prompted for:
- **User Name**: A friendly identifier (e.g., "prod-admin")
- **SSH Username**: The actual SSH username (e.g., "root")
- **SSH Password**: Your SSH password (hidden input)

Example:
```
User Name (identifier): prod-admin
SSH Username: root
SSH Password: ********
✓ User 'prod-admin' added successfully!
```

### Step 2: Add Your First Application

Applications define where logs are located on remote servers.

```bash
logx app add
```

You'll be prompted for:

1. **App Name**: Identifier for your application (e.g., "testapp")
2. **User**: Select from the users you've added
3. **Log file path**: Full path to log file (e.g., `/logs/testapp/testapp.log`)
4. **Log filename pattern**: Pattern with `{date}` placeholder
5. **Date format**: Go date format string
6. **Server IPs**: One or more server IPs

#### Example Configuration

```
App Name: testapp

Available users:
  1. prod-admin (username: root)
Select user (number): 1

Log file path: /logs/testapp/testapp.log

Log filename pattern with {date} placeholder: testapp-{date}.log

Date format examples:
  2006-01-02  -> 2025-09-10
  20060102    -> 20250910
  02-01-2006  -> 10-09-2025
Date format (Go format): 2006-01-02

Enter server IPs (one per line, empty line to finish):
Server IP: 192.168.0.1
Server IP: 192.168.0.2
Server IP:

✓ App 'testapp' added successfully!
```

### Step 3: Configure Editor (Optional)

By default, logx uses platform-specific editors. You can set a custom editor:

```bash
# Set VS Code
logx editor set code

# Set Notepad++
logx editor set notepad++

# Set Vim
logx editor set vim

# View current editor
logx editor show
```

## Using logx

### View Current Logs

View the current (non-dated) log file:

```bash
logx testapp
```

This will:
1. Connect to all configured servers
2. Download the log file from each server
3. Open each log file in your editor

### View Historical Logs

View logs for a specific date:

```bash
logx testapp 2025-09-10
```

This uses your configured date pattern and format to find the correct log file.

### View Logs from Specific Server

If you only want logs from one server:

```bash
# Current logs from specific server
logx testapp --server 192.168.0.1

# Historical logs from specific server
logx testapp 2025-09-10 --server 192.168.0.1
```

## Common Use Cases

### Scenario 1: Multiple Apps, Same User

If you have multiple applications using the same SSH credentials:

```bash
# Add user once
logx user add
# Name: prod-admin, Username: root

# Add multiple apps
logx app add  # App: webapp
logx app add  # App: apiservice
logx app add  # App: database

# View logs
logx webapp
logx apiservice 2025-09-10
logx database --server 10.0.0.5
```

### Scenario 2: Different Users for Different Environments

```bash
# Add users
logx user add  # Name: prod-admin
logx user add  # Name: dev-admin

# Add apps with different users
logx app add  # App: prod-app, User: prod-admin
logx app add  # App: dev-app, User: dev-admin
```

### Scenario 3: Different Date Formats

```bash
# App 1: testapp-2025-09-10.log
logx app add
# Pattern: testapp-{date}.log
# Format: 2006-01-02

# App 2: webapp.log-20250910
logx app add
# Pattern: webapp.log-{date}
# Format: 20060102

# App 3: api_2025_09_10.log
logx app add
# Pattern: api_{date}.log
# Format: 2006_01_02
```

## Troubleshooting

### Problem: "User not found"

```bash
# List configured users
logx user list

# If none exist, add one
logx user add
```

### Problem: "App not found"

```bash
# List configured apps
logx app list

# Add the missing app
logx app add
```

### Problem: "Failed to connect to server"

Check:
1. Server IP is correct
2. SSH service is running on port 22
3. Username and password are correct
4. Firewall allows SSH connections

```bash
# Test SSH manually
ssh username@192.168.0.1

# If that works, verify credentials in logx
logx user list
```

### Problem: "Log file not found"

Check:
1. Log path is correct
2. Date format matches actual log filename
3. Pattern includes `{date}` placeholder
4. User has read permissions

```bash
# SSH to server and check
ssh username@192.168.0.1
ls -l /logs/testapp/
```

### Problem: "Editor not opening"

```bash
# Set a specific editor
logx editor set code

# Or try platform defaults
logx editor set notepad++  # Windows
logx editor set gedit       # Linux
logx editor set open        # macOS
```

## Advanced Configuration

### Manual XML Editing

For advanced users, you can manually edit the configuration file:

**Location:**
- Linux/macOS: `~/.config/logx/config.xml`
- Windows: `%APPDATA%\logx\config.xml`

See `config.example.xml` for reference.

### Backing Up Configuration

```bash
# Linux/macOS
cp ~/.config/logx/config.xml ~/logx-backup.xml

# Windows
copy %APPDATA%\logx\config.xml %USERPROFILE%\logx-backup.xml
```

Note: Credentials are stored separately in the system keyring and cannot be exported.

## Security Best Practices

1. **Use Strong Passwords**: Ensure SSH passwords are strong
2. **Limit Access**: Only add servers you need access to
3. **Regular Updates**: Keep the tool updated
4. **SSH Keys**: Consider using SSH keys instead of passwords (future feature)
5. **Review Config**: Regularly review `logx app list` and `logx user list`

## Next Steps

- Add more applications: `logx app add`
- Add more users: `logx user add`
- Customize editor: `logx editor set <your-editor>`
- View logs: `logx <appname> [date]`

## Getting Help

```bash
# View all commands
logx help

# View usage
logx

# Check version
logx version
```

For more information, see [README.md](README.md).