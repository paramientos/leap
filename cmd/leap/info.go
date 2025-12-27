package main

import (
	"fmt"
	"strings"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [name]",
	Short: "Show detailed information about a connection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		conn, ok := cfg.Connections[name]
		if !ok {
			fmt.Printf("\n❌ Connection \033[1;36m%s\033[0m not found.\n\n", name)
			return
		}

		fmt.Println("\n⚡ \033[1;32mConnection Details\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")

		fmt.Printf("\n  \033[1m%-15s\033[0m \033[1;36m%s\033[0m", "Name:", conn.Name)
		if conn.Favorite {
			fmt.Print(" ⭐")
		}
		fmt.Println()

		fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Host:", conn.Host)
		fmt.Printf("  \033[1m%-15s\033[0m %s\n", "User:", conn.User)
		fmt.Printf("  \033[1m%-15s\033[0m %d\n", "Port:", conn.Port)

		if conn.IdentityFile != "" {
			fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Identity File:", conn.IdentityFile)
		}

		if len(conn.Tags) > 0 {
			fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Tags:", strings.Join(conn.Tags, ", "))
		}

		if conn.JumpHost != "" {
			fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Jump Host:", conn.JumpHost)
		}

		if len(conn.Tunnels) > 0 {
			fmt.Printf("  \033[1m%-15s\033[0m %d tunnels configured\n", "Tunnels:", len(conn.Tunnels))
		}

		fmt.Println("\n  \033[1m--- Stats ---\033[0m")
		fmt.Printf("  \033[1m%-15s\033[0m %d\n", "Usage Count:", conn.UsageCount)
		if !conn.LastUsed.IsZero() {
			fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Last Used:", conn.LastUsed.Format("2006-01-02 15:04:05"))
		}
		if !conn.CreatedAt.IsZero() {
			fmt.Printf("  \033[1m%-15s\033[0m %s\n", "Created At:", conn.CreatedAt.Format("2006-01-02 15:04:05"))
		}

		if conn.Notes != "" {
			fmt.Println("\n  \033[1m--- Notes ---\033[0m")
			fmt.Printf("  %s\n", conn.Notes)
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
