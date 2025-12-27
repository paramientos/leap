package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/paramientos/leap/internal/config"
	"github.com/paramientos/leap/internal/ssh"
	"github.com/paramientos/leap/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leap",
	Short: "⚡ LEAP - Modern SSH Connection Manager",
	Long: `
⚡ LEAP SSH MANAGER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

A modern CLI tool to manage your SSH connections with tags, 
fuzzy search, and an intuitive terminal interface.

Features:
  • 🔐 Secure encrypted configuration
  • 🏷️  Tag-based organization
  • 🔍 Fuzzy search & filtering
  • 🎨 Beautiful terminal UI
  • 🔀 Jump host support
  • 🚇 SSH tunnel management
`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		if len(args) > 0 {
			// Check if it's a connection name
			name := strings.Join(args, " ")
			if conn, ok := cfg.Connections[name]; ok {
				fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", name)
				ssh.Connect(conn, false)
				return
			}

			// Try partial match or tag match
			for _, conn := range cfg.Connections {
				if strings.Contains(strings.ToLower(conn.Name), strings.ToLower(name)) {
					fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", conn.Name)
					ssh.Connect(conn, false)
					return
				}
				for _, tag := range conn.Tags {
					if strings.EqualFold(tag, name) {
						fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", conn.Name)
						ssh.Connect(conn, false)
						return
					}
				}
			}
		}

		// Run TUI
		choice, err := tui.Run(cfg)
		if err != nil {
			fmt.Printf("\n❌ Error running TUI: %v\n\n", err)
			return
		}

		if choice != nil {
			err = ssh.Connect(*choice, false)
			if err != nil {
				fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
			}
		}
	},
}

func GetPassphrase() string {
	return os.Getenv("LEAP_MASTER_PASSWORD")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("⚡ LEAP SSH Manager v{{.Version}}\n")
}
