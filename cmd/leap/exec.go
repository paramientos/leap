package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [name] [command]",
	Short: "Execute a command on remote server(s)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		tag, _ := cmd.Flags().GetString("tag")

		var connsToExec []config.Connection
		var command string

		if all {
			command = strings.Join(args, " ")
			for _, conn := range cfg.Connections {
				connsToExec = append(connsToExec, conn)
			}
		} else if tag != "" {
			command = strings.Join(args, " ")
			for _, conn := range cfg.Connections {
				for _, t := range conn.Tags {
					if t == tag {
						connsToExec = append(connsToExec, conn)
						break
					}
				}
			}
		} else {
			name := args[0]
			command = strings.Join(args[1:], " ")

			if conn, ok := cfg.Connections[name]; ok {
				connsToExec = append(connsToExec, conn)
			} else {
				fmt.Printf("\n❌ Connection '\033[1;36m%s\033[0m' not found.\n", name)
				fmt.Println("\033[90mTip: Use 'leap list' to see all available connections\033[0m\n")
				return
			}
		}

		if len(connsToExec) == 0 {
			fmt.Println("\n\033[90mNo connections to execute on\033[0m\n")
			return
		}

		fmt.Println("\n⚡ \033[1;32mRemote Command Execution\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[90mCommand:\033[0m \033[1;35m%s\033[0m\n\n", command)

		for _, conn := range connsToExec {
			executeRemoteCommand(conn, command)
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	},
}

func executeRemoteCommand(conn config.Connection, command string) {
	fmt.Printf("\033[1;36m%s\033[0m (\033[33m%s\033[0m@\033[32m%s\033[0m)\n", conn.Name, conn.User, conn.Host)

	sshArgs := []string{}

	if conn.IdentityFile != "" {
		sshArgs = append(sshArgs, "-i", conn.IdentityFile)
	}

	sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", conn.Port))

	if conn.JumpHost != "" {
		sshArgs = append(sshArgs, "-J", conn.JumpHost)
	}

	sshArgs = append(sshArgs,
		"-o", "ConnectTimeout=10",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", conn.User, conn.Host),
		command,
	)

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	err := sshCmd.Run()
	if err != nil {
		fmt.Printf("\033[31m✗\033[0m Command failed: %v\n", err)
	} else {
		fmt.Printf("\033[32m✓\033[0m Command completed successfully\n")
	}
	fmt.Println()
}

func init() {
	execCmd.Flags().BoolP("all", "a", false, "Execute on all connections")
	execCmd.Flags().StringP("tag", "t", "", "Execute on connections with specific tag")
	rootCmd.AddCommand(execCmd)
}
