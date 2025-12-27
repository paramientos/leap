package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var scpCmd = &cobra.Command{
	Use:   "scp [connection] [local_path] [remote_path]",
	Short: "Transfer files using Leap connections",
	Long: `Transfer files between local and remote using Leap connection settings.
Example: leap scp myserver ./file.txt /tmp/`,
	Args: cobra.MinimumNArgs(1),
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

		if len(args) < 3 {
			fmt.Println("\n📂 Interactive File Manager coming soon!")
			fmt.Println("\033[90mUsage for now: leap scp [name] [local_src] [remote_dest]\033[0m\n")
			return
		}

		src := args[1]
		dest := args[2]

		fmt.Printf("\n🚀 Transferring \033[1;33m%s\033[0m to \033[1;36m%s:%s\033[0m...\n", src, name, dest)

		scpArgs := []string{"-P", fmt.Sprintf("%d", conn.Port)}
		if conn.IdentityFile != "" {
			scpArgs = append(scpArgs, "-i", conn.IdentityFile)
		}
		if conn.JumpHost != "" {
			scpArgs = append(scpArgs, "-J", conn.JumpHost)
		}

		remoteTarget := fmt.Sprintf("%s@%s:%s", conn.User, conn.Host, dest)
		scpArgs = append(scpArgs, src, remoteTarget)

		scpProcess := exec.Command("scp", scpArgs...)
		scpProcess.Stdout = os.Stdout
		scpProcess.Stderr = os.Stderr

		err = scpProcess.Run()
		if err != nil {
			fmt.Printf("\n❌ Transfer failed: %v\n\n", err)
		} else {
			fmt.Printf("\n\033[32m✓\033[0m File successfully transferred!\n\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(scpCmd)
}
