# logx - Interactive Remote Log Viewer 🚀

A beautiful command-line tool with an interactive TUI for viewing and managing log files from remote servers via SSH. Built with Go and Bubble Tea.

## ✨ Features

### 🎨 Interactive Terminal UI
- **Colorful menu navigation** with vim-style keybindings
- **Internal log viewer** - view logs directly in your terminal
- **Real-time search** with highlighted results
- **Save logs locally** with a single keypress
- **Intuitive navigation** using arrow keys or vim keys (j/k)

### 🔐 Security
- **System keyring integration** (Windows Credential Manager, macOS Keychain, Linux Secret Service)
- **Secure credential storage** - passwords never stored in plain text
- **SSH encrypted connections**

### 📋 Log Management
- **Multi-server support** - view logs from multiple servers simultaneously
- **Date-based log rolling** - support for various date formats in filenames
- **Pattern matching** - flexible log file naming patterns
- **Search functionality** - find text across log files instantly

### ⚙️ Configuration
- **User management** - add, list, delete SSH users
- **App management** - configure applications and their log locations
- **Custom editors** - use your preferred text editor
- **XML configuration** - easy-to-read and edit

## 📦 Installation

### Prerequisites
- Go 1.21 or higher
- SSH access to remote servers

### Quick Install

```bash
# Clone the repository
git clone https://github.com/jat-sandaruwan/logx.git
cd logx

# Build
make build

# Install system-wide (optional)
sudo make install
```

### Or use install script

```bash
chmod +x install.sh
./install.sh
```

## 🚀 Quick Start

### Launch Interactive TUI

Simply run:
```bash
logx
```

This launches the beautiful interactive menu where you can:
1. **Manage Users** - Add/view/delete SSH credentials
2. **Manage Apps** - Configure applications and log locations
3. **View Logs** - Browse and view logs interactively
4. **Settings** - Configure editor and preferences

### Command-Line Mode

You can also use traditional CLI commands:

```bash
# User management
logx user add
logx user list
logx user delete <username>

# App management
logx app add
logx app list
logx app update <appname>
logx app delete <appname>

# Editor configuration
logx editor set code
logx editor show

# Show version
logx version

# Show help
logx help
```

## 🎮 Interactive TUI Usage

### Main Menu Navigation

```
 ██╗      ██████╗  ██████╗ ██╗  ██╗
 ██║     ██╔═══██╗██╔════╝ ╚██╗██╔╝
 ██║     ██║   ██║██║  ███╗ ╚███╔╝ 
 ██║     ██║   ██║██║   ██║ ██╔██╗ 
 ███████╗╚██████╔╝╚██████╔╝██╔╝ ██╗
 ╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝
    Remote Log Viewer v1.0.0

👥 2 Users  •  📱 3 Apps

  ▶ 👤 User Management
    📱 App Management
    📋 View Logs
    ⚙️  Settings
    ❌ Exit

↑/↓: Navigate • Enter: Select • q: Quit
```

**Controls:**
- `↑/↓` or `j/k` - Navigate menu items
- `Enter` - Select option
- `Esc` - Go back to previous menu
- `q` or `Ctrl+C` - Quit

### Log Viewer Features

Once you select an app and view logs, you get a powerful internal viewer:

```
 📋 server1.example.com - app-2025-10-04.log 

   1  [INFO] Application started
   2  [INFO] Connecting to database
▶  3  [ERROR] Connection timeout
   4  [INFO] Retrying connection
   5  [INFO] Connected successfully

 Line 3/1250 | Match 1/5 

/: Search | n/N: Next/Prev | s: Save | q: Quit
```

**Viewer Controls:**

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate line by line |
| `Ctrl+U` or `PgUp` | Page up |
| `Ctrl+D` or `PgDn` | Page down |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `/` | Enter search mode |
| `n` | Next search result |
| `N` | Previous search result |
| `s` | Save log to local file |
| `q` or `Ctrl+C` | Close viewer |

### Search Feature

1. Press `/` to enter search mode
2. Type your search query
3. Press `Enter` to search
4. Results are highlighted in yellow
5. Use `n`/`N` to navigate between matches

### Saving Logs

Press `s` while viewing a log to save it locally. The file will be saved as:
```
server1.example.com_logs_app_app-2025-10-04.log.log
```

## 📝 Configuration

### Adding a User

**Via TUI:**
1. Launch `logx`
2. Select "User Management"
3. Select "Add User"
4. Follow the prompts

**Via CLI:**
```bash
logx user add
# Enter: User Name, SSH Username, SSH Password
```

### Adding an Application

**Via TUI:**
1. Launch `logx`
2. Select "App Management"
3. Select "Add App"
4. Follow the interactive prompts

**Via CLI:**
```bash
logx app add
```

You'll configure:
- **App Name:** Identifier (e.g., `myapp`)
- **User:** SSH user to use
- **Log Path:** Base path (e.g., `/var/log/myapp/app.log`)
- **Pattern:** Filename pattern with `{date}` (e.g., `app-{date}.log`)
- **Date Format:** Go date format (e.g., `2006-01-02`)
- **Servers:** IP addresses (one per line)

### Configuration File

Located at:
- **Linux/macOS:** `~/.config/logx/config.xml`
- **Windows:** `%APPDATA%\logx\config.xml`

Example:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<config>
  <users>
    <user id="prod" name="prod" username="root"/>
  </users>
  <apps>
    <app name="webapp">
      <user-ref>prod</user-ref>
      <log-path>/var/log/webapp/app.log</log-path>
      <log-pattern>app-{date}.log</log-pattern>
      <date-format>2006-01-02</date-format>
      <servers>
        <server>192.168.1.10</server>
        <server>192.168.1.11</server>
      </servers>
    </app>
  </apps>
  <editor>code</editor>
</config>
```

## 🎯 Use Cases

### Scenario 1: View Today's Logs

```bash
logx
# Select "View Logs"
# Choose your app
# Select "All Servers" or specific server
# Press Enter (defaults to today)
```

### Scenario 2: Search for Errors

1. View a log file
2. Press `/`
3. Type `ERROR`
4. Press `Enter`
5. Use `n` to navigate through all errors

### Scenario 3: Compare Logs from Multiple Servers

1. View logs from Server 1
2. Press `s` to save
3. Press `q` to exit
4. View logs from Server 2
5. Press `s` to save
6. Open both files in your editor for comparison

### Scenario 4: Historical Log Analysis

```bash
logx
# Select "View Logs"
# Choose your app
# Choose server
# Enter date: 2025-09-15
# View, search, and save as needed
```

## 🎨 Customization

### Custom Editor

**Via TUI:**
1. Settings → Set Editor
2. Enter editor command

**Via CLI:**
```bash
logx editor set "code -w"
logx editor set "vim"
logx editor set "notepad++"
```

### Color Scheme

The TUI uses a carefully chosen color palette:
- **Primary:** Purple (`#7D56F4`)
- **Success:** Green (`#04B575`)
- **Warning:** Orange (`#FFA500`)
- **Error:** Red (`#FF0000`)
- **Highlight:** Yellow (for search results)

## 🔧 Troubleshooting

### TUI Not Displaying Correctly

```bash
# Ensure terminal supports colors
echo $TERM
# Should be: xterm-256color or similar

# If not, set it:
export TERM=xterm-256color
```

### SSH Connection Issues

1. Test SSH manually first:
   ```bash
   ssh username@server
   ```

2. Check firewall rules
3. Verify credentials in keyring

### Log File Not Found

1. Verify the log path exists on the server
2. Check date format matches actual files
3. Ensure pattern contains `{date}` placeholder
4. Test with `logx app list` to view configuration

### Search Not Working

- Search is case-insensitive
- Ensure you press `Enter` after typing query
- Try simpler search terms

## 📊 Date Format Reference

| Format | Output | Use Case |
|--------|--------|----------|
| `2006-01-02` | 2025-10-04 | ISO 8601 |
| `20060102` | 20251004 | Compact |
| `02-01-2006` | 04-10-2025 | DD-MM-YYYY |
| `01/02/2006` | 10/04/2025 | US format |
| `2006_01_02` | 2025_10_04 | Underscore |
| `2006-01` | 2025-10 | Monthly logs |
| `2006` | 2025 | Yearly logs |

## 🚀 Advanced Features

### Parallel Viewing (Coming Soon)
- View multiple log files side-by-side
- Real-time log tailing
- Log diff mode

### Filtering (Coming Soon)
- Filter by log level (INFO, WARN, ERROR)
- Time range filtering
- Custom regex filters

### Export Options (Coming Soon)
- Export search results
- Generate reports
- Email/Slack notifications

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
git clone https://github.com/jat-sandaruwan/logx.git
cd logx
go mod download
go build -o logx cmd/logx/main.go
```

### Adding New Features

Key files:
- `internal/ui/*.go` - TUI components
- `internal/viewer/viewer.go` - Log viewer
- `cmd/logx/main.go` - CLI entry point

## 📄 License

MIT License - see LICENSE.md

## 👤 Author

**Thilina Sandaruwan**

## 🌟 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [go-keyring](https://github.com/zalando/go-keyring) - Secure credential storage

## 📞 Support

- **Issues:** https://github.com/jat-sandaruwan/logx/issues
- **Discussions:** https://github.com/jat-sandaruwan/logx/discussions

## 🎉 Star the Project

If you find logx useful, please give it a ⭐️ on GitHub!

---

**Made with ❤️ for developers who love beautiful CLIs**