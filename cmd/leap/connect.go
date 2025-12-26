package main

import (
	"fmt"
	"strings"

	"github.com/paramientos/leap/internal/config"
	"github.com/paramientos/leap/internal/ssh"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [name]",
	Short: "Connect to a saved SSH host",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a connection name or use 'sshm list'")
			return
		}

		name := args[0]
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		conn, ok := cfg.Connections[name]
		if !ok {
			// Try fuzzy match or look for similar names
			fmt.Printf("Connection '%s' not found. Searching...\n", name)
			var found bool
			for k, v := range cfg.Connections {
				if strings.Contains(k, name) {
					fmt.Printf("Found match: %s. Connecting...\n", k)
					conn = v
					found = true
					break
				}
			}
			if !found {
				fmt.Println("No matching connection found.")
				return
			}
		}

		err = ssh.Connect(conn)
		if err != nil {
			fmt.Printf("SSH Connection closed with error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
