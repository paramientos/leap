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
			fmt.Println("\n❌ Please specify a connection name or use 'leap list'")
			fmt.Println("\033[90mUsage: leap connect [name]\033[0m\n")
			return
		}

		name := args[0]
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		conn, ok := cfg.Connections[name]
		if !ok {
			// Try fuzzy match or look for similar names
			fmt.Printf("\n🔍 Connection '\033[1;36m%s\033[0m' not found. Searching...\n", name)
			var found bool
			for k, v := range cfg.Connections {
				if strings.Contains(k, name) {
					fmt.Printf("\033[32m✓\033[0m Found match: \033[1;36m%s\033[0m. Connecting...\n\n", k)
					conn = v
					found = true
					break
				}
			}
			if !found {
				fmt.Println("\n❌ No matching connection found.")
				fmt.Println("\033[90mTip: Use 'leap list' to see all available connections\033[0m\n")
				return
			}
		} else {
			fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m (\033[33m%s\033[0m@\033[32m%s\033[0m)...\n\n", name, conn.User, conn.Host)
		}

		cfg.UpdateLastUsed(name)
		config.SaveConfig(cfg, GetPassphrase())

		record, _ := cmd.Flags().GetBool("record")

		err = ssh.Connect(conn, record)
		if err != nil {
			fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
		}
	},
}

func init() {
	connectCmd.Flags().BoolP("record", "r", false, "Record session")
	rootCmd.AddCommand(connectCmd)
}
