# logx Project Summary

## Overview

**logx** is a command-line tool written in Go that allows you to view and manage log files from remote servers via SSH. It provides a beautiful terminal UI using Bubble Tea and securely stores credentials in the system keyring.

## Project Structure

```
logx/
├── cmd/
│   └── logx/
│       └── main.go              # CLI entry point, command routing
├── internal/
│   ├── config/
│   │   └── config.go            # XML configuration management
│   ├── editor/
│   │   └── editor.go            # Platform-specific editor handling
│   ├── ssh/
│   │   └── ssh.go               # SSH client, file operations
│   ├── ui/
│   │   ├── app.go               # App management UI
│   │   └── user.go              # User management UI
│   ├── vault/
│   │   └── vault.go             # System keyring integration
│   └── viewer/
│       └── viewer.go            # Log viewing and downloading
├── build/                        # Build artifacts (created on build)
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── Makefile                      # Build automation
├── .gitignore                    # Git ignore rules
├── config.example.xml            # Example configuration
├── install.sh                    # Unix installation script
├── install.ps1                   # Windows installation script
├── README.md                     # Main documentation
└── SETUP.md                      # Setup guide

```

## Core Components

### 1. Configuration Management (`internal/config/`)

- **Purpose**: Manages XML-based configuration for users and apps
- **Key Features**:
  - Platform-specific config paths (~/.config/logx/ or %APPDATA%\logx\)
  - User management (add, list, delete, get)
  - App management (add, list, update, delete, get)
  - Editor preference storage

### 2. Credential Vault (`internal/vault/`)

- **Purpose**: Secure credential storage using system keyring
- **Platforms**:
  - Windows: Windows Credential Manager
  - macOS: Keychain
  - Linux: Secret Service API
- **Operations**: Store, Get, Delete, Exists

### 3. SSH Client (`internal/ssh/`)

- **Purpose**: Handle SSH connections and file operations
- **Features**:
  - Password authentication
  - File existence checking
  - File downloading
  - Pattern-based file listing
- **Implementation**: Uses `golang.org/x/crypto/ssh`

### 4. Editor Management (`internal/editor/`)

- **Purpose**: Open files in platform-appropriate editors
- **Default Editors**:
  - Windows: Notepad++ → Notepad
  - macOS: VS Code → Sublime → TextEdit
  - Linux: VS Code → gedit → kate → nano → vim
- **Custom**: Supports user-defined editor commands

### 5. UI Components (`internal/ui/`)

- **Purpose**: Interactive terminal interfaces
- **Technologies**: Bubble Tea, Lipgloss
- **Features**:
  - Interactive user addition
  - Interactive app addition/update
  - Listing utilities

### 6. Log Viewer (`internal/viewer/`)

- **Purpose**: Core log viewing functionality
- **Process**:
  1. Load app configuration
  2. Get user credentials from vault
  3. Connect to server(s) via SSH
  4. Download log file(s) to temp directory
  5. Open in configured editor
- **Features**:
  - Date-based log file resolution
  - Multi-server support
  - Server filtering

### 7. Main CLI (`cmd/logx/`)

- **Purpose**: Command-line interface and routing
- **Commands**:
  - `user add|list|delete` - User management
  - `app add|list|update|delete` - App management
  - `editor set|show` - Editor configuration
  - `<appname> [date] [--server IP]` - View logs
  - `help` - Show help
  - `version` - Show version

## Data Flow

### Adding a User
```
User Input → UI Prompt → Config (XML) + Vault (Keyring)
```

### Adding an App
```
User Input → UI Prompt → Load Users → Save to Config (XML)
```

### Viewing Logs
```
CLI → Load Config → Get Credentials from Vault
    → SSH Connect → Check File → Download
    → Open in Editor
```

## Configuration Format

### XML Structure
```xml
<config>
  <users>
    <user id="..." name="..." username="..."/>
  </users>
  <apps>
    <app name="...">
      <user-ref>...</user-ref>
      <log-path>...</log-path>
      <log-pattern>...-{date}.log</log-pattern>
      <date-format>2006-01-02</date-format>
      <servers>
        <server>IP</server>
      </servers>
    </app>
  </apps>
  <editor>command</editor>
</config>
```

### Log Pattern System

- **Pattern**: Contains `{date}` placeholder (e.g., `app-{date}.log`)
- **Date Format**: Go time format (e.g., `2006-01-02`)
- **Resolution**: Pattern + Format + Date → Filename

Examples:
- Pattern: `app-{date}.log`, Format: `2006-01-02`, Date: 2025-09-10 → `app-2025-09-10.log`
- Pattern: `app.log-{date}`, Format: `20060102`, Date: 2025-09-10 → `app.log-20250910`

## Security Features

1. **Credential Storage**: System keyring (not plain text)
2. **Config Permissions**: 0600 (user read/write only)
3. **SSH**: Encrypted connections
4. **Temp Files**: Downloaded logs in system temp directory

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| `github.com/zalando/go-keyring` | System keyring access |
| `golang.org/x/crypto/ssh` | SSH client |
| `golang.org/x/term` | Terminal utilities |

## Build System

### Makefile Targets
- `make build` - Build for current platform
- `make build-all` - Build for all platforms
- `make install` - Install to /usr/local/bin
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make deps` - Download dependencies

### Cross-Platform Builds
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o logx-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o logx.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o logx-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o logx-darwin-arm64
```

## Usage Patterns

### Basic Workflow
1. **Setup**: `logx user add` → `logx app add`
2. **View**: `logx <appname>` or `logx <appname> 2025-09-10`
3. **Manage**: `logx app list`, `logx user list`

### Advanced Features
- Server filtering: `--server 192.168.0.1`
- Custom editor: `logx editor set code`
- Updates: `logx app update <name>`

## Extension Points

### Future Enhancements
1. **SSH Keys**: Support SSH key authentication
2. **Compression**: Handle compressed log files (.gz, .zip)
3. **Search**: Built-in log searching
4. **Tail**: Real-time log tailing
5. **Filters**: Pattern-based filtering
6. **Export**: Export logs to local directory
7. **Diff**: Compare logs from multiple servers
8. **Alerts**: Email/Slack notifications
9. **Web UI**: Optional web interface
10. **Plugins**: Plugin system for custom processors

### Code Extension Points
- `internal/ssh/ssh.go`: Add new SSH operations (keys, compression)
- `internal/viewer/viewer.go`: Add log processing (search, filter, merge)
- `internal/editor/editor.go`: Add new editor support
- `internal/config/config.go`: Extend configuration options
- `cmd/logx/main.go`: Add new commands

## Error Handling Strategy

### Connection Errors
- Graceful failure per server
- Continue with remaining servers
- Clear error messages to user

### Configuration Errors
- Validation on add/update
- Helpful error messages
- Suggest corrections

### File Errors
- Check existence before download
- Handle permission issues
- Report which server failed

## Testing Strategy

### Unit Tests (To Implement)
```go
// internal/config/config_test.go
TestLoadConfig()
TestSaveConfig()
TestAddUser()
TestGetApp()

// internal/vault/vault_test.go
TestStoreCredentials()
TestGetCredentials()

// internal/ssh/ssh_test.go (with mock)
TestConnect()
TestDownloadFile()
```

### Integration Tests
- Mock SSH server for testing
- Test full workflow: add user → add app → view logs
- Test error conditions

### Manual Testing Checklist
- [ ] Add user with valid credentials
- [ ] Add user with invalid credentials
- [ ] Add app with all required fields
- [ ] View current logs
- [ ] View historical logs
- [ ] View logs from specific server
- [ ] Update app configuration
- [ ] Delete user (check vault cleanup)
- [ ] Delete app
- [ ] Set custom editor
- [ ] Test on Windows
- [ ] Test on macOS
- [ ] Test on Linux

## Performance Considerations

### Current Implementation
- **Serial Downloads**: Downloads logs from servers sequentially
- **Temp Storage**: Uses system temp directory
- **Memory**: Streams file downloads (not loaded into memory)

### Optimization Opportunities
1. **Parallel Downloads**: Use goroutines for concurrent server connections
2. **Caching**: Cache recently accessed logs
3. **Compression**: Support downloading compressed logs
4. **Streaming**: Stream large files directly to editor

## Deployment

### Binary Distribution
```bash
# Create releases for all platforms
make build-all

# Package
tar -czf logx-linux-amd64.tar.gz build/logx-linux-amd64
zip logx-windows-amd64.zip build/logx-windows-amd64.exe
tar -czf logx-darwin-amd64.tar.gz build/logx-darwin-amd64
```

### GitHub Release
1. Tag version: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. Create release on GitHub
4. Upload binaries

### Package Managers (Future)
- **Homebrew**: Create formula
- **Chocolatey**: Create package
- **apt/yum**: Create .deb/.rpm packages

## Troubleshooting Guide

### Common Issues

#### 1. Keyring Access Issues
**Linux**: Ensure gnome-keyring or similar is running
```bash
# Check if keyring daemon is running
ps aux | grep keyring

# Start gnome-keyring if needed
gnome-keyring-daemon --start
```

#### 2. SSH Connection Failures
- Check firewall rules
- Verify SSH service is running
- Test with standard SSH client first
- Check for IP whitelisting

#### 3. Config File Permissions
```bash
# Fix permissions if needed
chmod 600 ~/.config/logx/config.xml
```

#### 4. Editor Not Found
```bash
# Check editor installation
which code
which notepad++
which vim

# Set full path if needed
logx editor set /usr/bin/vim
```

## Development Workflow

### Setup Development Environment
```bash
# Clone repository
git clone https://github.com/jatsandaruwan/logx.git
cd logx

# Install dependencies
go mod download

# Build
go build -o logx cmd/logx/main.go

# Run
./logx version
```

### Making Changes
1. Create feature branch: `git checkout -b feature/name`
2. Make changes
3. Test locally
4. Commit: `git commit -m "Description"`
5. Push: `git push origin feature/name`
6. Create Pull Request

### Code Style
- Follow Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

## Architecture Decisions

### Why XML for Configuration?
- Human-readable and editable
- Easy to validate
- Standard library support
- Hierarchical structure fits use case

### Why System Keyring?
- Platform-native security
- Better than plain text or encrypted files
- User familiar with OS credential management
- No custom encryption needed

### Why Bubble Tea?
- Modern, elegant TUI framework
- Active development
- Good documentation
- Makes CLI tools more user-friendly

### Why Not SSH Keys Initially?
- Password auth is simpler for first version
- Many users have passwords available
- SSH keys can be added in v2

### Why Download Files Instead of Streaming?
- Editors need local files
- Allows offline viewing
- Can open multiple files simultaneously
- Simpler implementation

## Metrics & Monitoring (Future)

### Potential Metrics
- Number of log views per day
- Average download time
- Failed connection rate
- Most viewed applications
- Storage usage

### Implementation Ideas
- Optional telemetry (opt-in)
- Local log file for usage stats
- Export to monitoring systems

## Documentation

### Provided Documents
1. **README.md** - Main documentation
2. **SETUP.md** - Setup guide for new users
3. **PROJECT_SUMMARY.md** - This file (architecture overview)
4. **config.example.xml** - Configuration example

### Additional Documentation Needed
1. **CONTRIBUTING.md** - Contribution guidelines
2. **CHANGELOG.md** - Version history
3. **API.md** - Internal API documentation
4. **SECURITY.md** - Security policy

## Version History

### v1.0.0 (Initial Release)
- User management (add, list, delete)
- App management (add, list, update, delete)
- SSH password authentication
- System keyring integration
- Multi-server support
- Date-based log file patterns
- Platform-specific editor support
- Custom editor configuration
- Interactive CLI with Bubble Tea

### Future Versions

#### v1.1.0 (Planned)
- SSH key authentication
- Parallel server downloads
- Log file compression support
- Search within logs
- Config import/export

#### v1.2.0 (Planned)
- Real-time log tailing
- Log filtering and parsing
- Web interface (optional)
- Plugin system

#### v2.0.0 (Planned)
- Log aggregation and merging
- Cloud storage integration
- Team sharing features
- Advanced search and analytics

## Resources

### Documentation
- Go Documentation: https://golang.org/doc/
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- SSH Package: https://pkg.go.dev/golang.org/x/crypto/ssh
- Keyring: https://github.com/zalando/go-keyring

### Similar Projects
- `ssh-tail`: Tail logs over SSH
- `multitail`: Multiple log file viewer
- `lnav`: Log file navigator

## Contact & Support

- **Author**: Thilina Sandaruwan
- **Repository**: https://github.com/jat-sandaruwan/logx
- **Issues**: https://github.com/jat-sandaruwan/logx/issues

## License

MIT License - See LICENSE file for details