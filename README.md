# ⚡ LEAP SSH Manager

A modern, beautiful CLI tool to manage your SSH connections with an intuitive terminal interface inspired by Laravel's elegant design.

![LEAP SSH Manager](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## ✨ Features

- 🔐 **Secure encrypted configuration** - Your main config is safely encrypted
- 🏷️ **Tag-based organization** - Organize connections with custom tags
- 🔍 **Fuzzy search & filtering** - Find connections quickly
- 🎨 **Beautiful terminal UI** - Modern, colorful interface inspired by Laravel
- 🔀 **Jump host support** - Connect through bastion hosts
- 🚇 **SSH tunnel management** - Create and manage SSH tunnels easily
- ⚡ **Fast and lightweight** - Built with Go for maximum performance
- ⭐ **Favorites system** - Mark frequently used connections
- 📝 **Connection notes** - Add notes to your connections
- 🧪 **Health checks** - Test connections and measure latency
- 📤 **Plain-text Export/Import** - Easily backup and share configurations
- 🖥️ **Remote execution** - Run commands on multiple servers
- 📁 **File transfer** - Upload/download files via SCP
- ✏️ **Edit connections** - Update existing connections easily
- 🗑️ **Bulk operations** - Delete multiple connections at once

## 📦 Installation

### From Source

```bash
git clone https://github.com/paramientos/leap.git
cd leap
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
- 🌐 Hostname
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

### Encryption

Connections are automatically encrypted using the [age](https://github.com/FiloSottile/age) encryption format if you have a master password set.

Set a master password using the `LEAP_MASTER_PASSWORD` environment variable. You can set it globally in your `~/.zshrc` or `~/.bashrc`:

```bash
export LEAP_MASTER_PASSWORD="your-secure-password"
```

Or provide it temporarily for a single command:

```bash
LEAP_MASTER_PASSWORD="your-secure-password" leap list
```

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
- Terminal with true color support

### Building

```bash
go build -o leap ./cmd/leap
```

### Running Tests

```bash
go test ./...
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
