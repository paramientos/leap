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
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		conn, ok := cfg.Connections[name]
		if !ok {
			fmt.Printf("Connection '%s' not found.\n", name)
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

		fmt.Printf("✓ Opening tunnel: %s -> %s\n", tunnelSpec, name)
		fmt.Println("Press Ctrl+C to close the tunnel.")

		c := exec.Command("ssh", sshArgs...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		err = c.Run()
		if err != nil {
			fmt.Printf("Tunnel closed: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
}
