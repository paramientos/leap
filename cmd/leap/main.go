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
	Short: "SSH Connection Manager",
	Long:  `A CLI tool to manage your SSH connections with tags, fuzzy search, and more.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if len(args) > 0 {
			// Check if it's a connection name
			name := strings.Join(args, " ")
			if conn, ok := cfg.Connections[name]; ok {
				ssh.Connect(conn)
				return
			}

			// Try partial match or tag match
			for _, conn := range cfg.Connections {
				if strings.Contains(strings.ToLower(conn.Name), strings.ToLower(name)) {
					ssh.Connect(conn)
					return
				}
				for _, tag := range conn.Tags {
					if strings.EqualFold(tag, name) {
						ssh.Connect(conn)
						return
					}
				}
			}
		}

		// Run TUI
		choice, err := tui.Run(cfg)
		if err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			return
		}

		if choice != nil {
			err = ssh.Connect(*choice)
			if err != nil {
				fmt.Printf("SSH Connection closed with error: %v\n", err)
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
	// Root flags if any
}
