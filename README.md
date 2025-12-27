# ⚡ LEAP SSH Manager

**The Ultimate SSH Connection Manager** - A modern, feature-rich CLI tool that goes beyond simple SSH management. Monitor servers in real-time, capture snapshots, share connections via QR codes, record sessions, and manage everything from a beautiful terminal UI.

Built with Go for maximum performance. Inspired by Laravel's elegant design philosophy.

![LEAP SSH Manager](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## 🎯 Why LEAP?

LEAP isn't just another SSH manager - it's a complete DevOps command center:
- 📸 **Snapshot & Compare** servers to track changes over time
- 📱 **QR Code Sharing** for instant connection distribution
- ⏺️ **Session Recording** for documentation and auditing
- 📊 **Live Monitoring** of server resources in beautiful TUIs
- 🔑 **Automated Key Management** with one-command setup

All in a single binary with zero dependencies!

## ✨ Features

- 🔐 **Secure encrypted configuration** - Your main config is safely encrypted
- 📸 **Server Snapshots** - Capture complete server state (OS, packages, services, ports)
- 📱 **QR Share** - Share connections via **QR Codes** ⚡
- ⏺️ **Session Recording** - Record and replay SSH sessions ⏺️
- 📊 **Real-time Monitoring** - Watch server Load, RAM and Uptime in a live TUI
- 🔑 **Self-Managed SSH Keys** - Generate and push Leap-specific SSH keys automatically
- 🏷️ **Tag-based & Group organization** - Organize connections with tags and folders
- 🔍 **Fuzzy search & filtering** - Find connections quickly
- 🎨 **Beautiful terminal UI** - Modern, colorful interface inspired by Laravel
- 🔀 **Jump host support** - Connect through bastion hosts
- 🚇 **SSH tunnel management** - Create and manage SSH tunnels easily
- 📂 **Smart SCP** - Transfer files using saved connection parameters
- 🧪 **Health checks** - Test connections and measure latency with visual bars
- 📤 **Plain-text Export/Import** - Easily backup and share configurations
- 📁 **SSH Config Import** - Migrate from `~/.ssh/config` in one command

## 📦 Installation

### Download Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/paramientos/leap/releases):

- **Linux (AMD64/ARM64)**
- **macOS (Intel/Apple Silicon)**
- **Windows (AMD64/ARM64)**

```bash
# Extract and install (Linux/macOS)
tar -xzf leap-*.tar.gz
sudo mv leap /usr/local/bin/

# Make executable
chmod +x /usr/local/bin/leap
```

### Build from Source

```bash
git clone https://github.com/paramientos/leap.git
cd leap

# Using Makefile (recommended)
make build          # Build for current platform
make install        # Build and install to /usr/local/bin
make build-all      # Build for all platforms
make release        # Create release archives

# Or using Go directly
go build -o leap ./cmd/leap
sudo mv leap /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/paramientos/leap/cmd/leap@latest
```

## 🚀 Quick Start

### Interactive TUI Mode

Simply run `leap` to launch the interactive terminal UI:

```bash
leap
```

### Add a New Connection

```bash
leap add myserver
```

You'll be prompted for:
- 🌐 Hostname/IP
- 👤 User
- 🔌 Port
- 🔐 Password (optional)
- 🔑 SSH Key Path (optional)
- 🏷️ Tags
- 🔀 Jump Host (optional)

### List All Connections

```bash
leap list
```

Filter by tag:

```bash
leap list --tag production
```

### Connect to a Server

```bash
leap connect myserver
```

Or use fuzzy matching:

```bash
leap myserver
```

### Edit a Connection

```bash
leap edit myserver
```

### Delete Connection(s)

```bash
leap delete oldserver
leap delete server1 server2 server3  # Multiple
leap delete myserver --force         # Skip confirmation
```

### Test Connections

```bash
leap test myserver              # Test single connection
leap test --all                 # Test all connections
leap test --tag production      # Test by tag
```

### Manage Favorites

```bash
leap favorite myserver          # Toggle favorite
leap favorites                  # List all favorites
```

### Add Notes

```bash
leap notes myserver             # View notes
leap notes myserver --edit      # Edit notes
```

### Remote Command Execution

```bash
leap exec myserver "uptime"
leap exec --all "df -h"
leap exec --tag web "systemctl status nginx"
```

### File Transfer

```bash
# Upload
leap upload myserver ./local.txt /remote/path/
leap upload myserver ./folder/ /remote/ --recursive

# Download
leap download myserver /remote/file.txt ./local/
leap download myserver /remote/folder/ ./ --recursive
```

### Export/Import

### Server Snapshots

Capture and compare complete server states for change tracking and auditing.

```bash
# Capture a snapshot
leap snapshot myserver -o snapshot.json

# Capture with installed packages (slower)
leap snapshot myserver -o snapshot.json --packages

# Compare two snapshots
leap diff snapshot1.json snapshot2.json

# YAML format
leap snapshot myserver -f yaml -o snapshot.yaml
```

### Share Connections

Share connection details via QR code or encrypted short-codes.

```bash
# Generate QR code and share code
leap share myserver

# Receiver imports with:
leap import-code [base64-code]
```

### Session Recording

Record and replay SSH sessions for documentation or auditing.

```bash
# Record a session
leap connect myserver --record

# List recordings
leap history

# Replay a session
leap replay myserver_20231227_153045
```

### File Transfer

Transfer files using your saved connection settings.

```bash
# Upload file
leap scp myserver ./local-file.txt /remote/path/

# The command automatically uses your saved port, keys, and jump hosts
```

### Health & Monitoring

Check if your servers are alive or watch their resources in real-time.

```bash
# Health check (with visual latency bars)
leap test --all

# Live Resource Monitor (CPU, RAM, Uptime)
leap monitor
leap monitor server1 server2
```

### SSH Key Wizard

Automate passwordless login by generating and pushing LEAP-specific keys.

```bash
# Generates a key if missing and pushes it to the server
leap push-key myserver
```

### Import from SSH Config

Migrate your existing connections from your system's SSH configuration.

```bash
leap import-ssh
```

### Export/Import

Backup or share your connections in JSON or YAML format. Note that items are exported in **plain-text** for easy sharing and manual editing.

```bash
# Export (Decrypted output)
leap export backup.json
leap export backup.yaml --format yaml

# Import (Reads plain-text)
leap import backup.json
leap import backup.yaml --merge  # Update existing
```

### Create SSH Tunnel

```bash
leap tunnel myserver 8080:localhost:80
```

## 🎨 Screenshots

### Main TUI Interface
The interactive terminal UI provides a beautiful, modern interface for managing your SSH connections:

- **Left Panel**: List of all your connections with fuzzy search
- **Right Panel**: Detailed information about the selected connection
- **Color-coded**: Easy to read with Laravel-inspired color scheme

### List Command
```
⚡ LEAP SSH MANAGER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

NAME          CONNECTION                    TAGS
────          ──────────                    ────
production    user@prod.example.com:22      #prod #web
staging       user@staging.example.com:22   #staging #web
database      admin@db.example.com:22       #prod #database

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓ Total connections: 3
```

## 🔧 Configuration

LEAP stores your connections in a configuration file at:
- macOS/Linux/Windows: `~/.leap/connections.yaml`

### Encryption (Optional)

By default, LEAP stores your connections in plain-text YAML. If you want to encrypt your configuration file, set a master password using the `LEAP_MASTER_PASSWORD` environment variable:

```bash
export LEAP_MASTER_PASSWORD="your-secure-password"
```

If set, your configuration will be automatically encrypted using the [age](https://github.com/FiloSottile/age) format.

## 📖 Usage Examples

### Quick Connect by Name
```bash
leap production
```

### Quick Connect by Tag
```bash
leap web
```

### Port Forwarding
```bash
# Forward local port 3306 to remote MySQL
leap tunnel database 3306:localhost:3306
```

### Using Jump Hosts
When adding a connection, specify a jump host:
```
🔀 Jump Host: bastion.example.com
```

## 🎯 Keyboard Shortcuts (TUI Mode)

- `↑/↓` or `j/k` - Navigate through connections
- `/` - Start filtering/searching
- `Enter` - Connect to selected server
- `q` or `Ctrl+C` - Quit

## 🛠️ Development

### Prerequisites

- Go 1.24 or higher
- Make (optional, but recommended)

### Building

```bash
# Using Makefile
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Install to /usr/local/bin

# Or using Go directly
go build -o leap ./cmd/leap
```

### Running Tests

```bash
make test
# or
go test ./...
```

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Build and install
make clean          # Clean build artifacts
make test           # Run tests
make deps           # Update dependencies
make release        # Create release archives
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by Laravel's beautiful CLI design
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Uses [Cobra](https://github.com/spf13/cobra) for CLI framework
- Encryption powered by [age](https://github.com/FiloSottile/age)

## 📧 Contact

Project Link: [https://github.com/paramientos/leap](https://github.com/paramientos/leap)

---

Made with ❤️ and ☕ by the LEAP team
