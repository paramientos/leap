package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var tunnelCmd = &cobra.Command{
	Use:   "tunnel [name] [local_port:remote_host:remote_port]",
	Short: "Open an SSH tunnel",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		tunnelSpec := args[1]

		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		conn, ok := cfg.Connections[name]
		if !ok {
			fmt.Printf("\n❌ Connection '\033[1;36m%s\033[0m' not found.\n", name)
			fmt.Println("\033[90mTip: Use 'leap list' to see all available connections\033[0m\n")
			return
		}

		sshArgs := []string{}
		if conn.IdentityFile != "" {
			sshArgs = append(sshArgs, "-i", conn.IdentityFile)
		}
		sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", conn.Port))
		if conn.JumpHost != "" {
			sshArgs = append(sshArgs, "-J", conn.JumpHost)
		}

		sshArgs = append(sshArgs, "-L", tunnelSpec)
		sshArgs = append(sshArgs, "-N") // Do not execute a remote command
		sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", conn.User, conn.Host))

		fmt.Println("\n⚡ \033[1;32mSSH Tunnel\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Opening tunnel: \033[1;35m%s\033[0m → \033[1;36m%s\033[0m\n", tunnelSpec, name)
		fmt.Printf("\033[90mPress Ctrl+C to close the tunnel\033[0m\n\n")

		c := exec.Command("ssh", sshArgs...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		err = c.Run()
		if err != nil {
			fmt.Printf("\n❌ Tunnel closed: %v\n\n", err)
		} else {
			fmt.Println("\n\033[32m✓\033[0m Tunnel closed successfully\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
}
