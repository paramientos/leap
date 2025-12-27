package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var pushKeyCmd = &cobra.Command{
	Use:   "push-key [name]",
	Short: "Upload your public key to a remote server",
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

		home, _ := os.UserHomeDir()
		leapKeyPath := filepath.Join(home, ".leap", "id_ed25519")
		leapPubKeyPath := leapKeyPath + ".pub"

		// 1. Generate key if it doesn't exist
		if _, err := os.Stat(leapKeyPath); os.IsNotExist(err) {
			fmt.Printf("\n🔑 No Leap SSH key found. Generating a new one at \033[1;36m%s\033[0m...\n", leapKeyPath)

			err := os.MkdirAll(filepath.Dir(leapKeyPath), 0700)
			if err != nil {
				fmt.Printf("\n❌ Failed to create directory: %v\n", err)
				return
			}

			genKey := exec.Command("ssh-keygen", "-t", "ed25519", "-f", leapKeyPath, "-N", "", "-C", "leap-ssh-manager")
			err = genKey.Run()
			if err != nil {
				fmt.Printf("\n❌ Failed to generate SSH key: %v\n", err)
				return
			}
			fmt.Println("✅ SSH key generated successfully.")
		}

		fmt.Printf("\n🚀 Pushing public key to \033[1;36m%s\033[0m...\n", name)

		sshCopyId := exec.Command("ssh-copy-id",
			"-i", leapPubKeyPath,
			"-p", fmt.Sprintf("%d", conn.Port),
			fmt.Sprintf("%s@%s", conn.User, conn.Host),
		)

		sshCopyId.Stdout = os.Stdout
		sshCopyId.Stderr = os.Stderr
		sshCopyId.Stdin = os.Stdin

		err = sshCopyId.Run()
		if err != nil {
			fmt.Printf("\n❌ Failed to push key: %v\n", err)
			fmt.Println("\033[90mNote: You might need to enter your password for the last time.\033[0m\n")
			return
		}

		// 2. Update connection to use this IdentityFile
		conn.IdentityFile = leapKeyPath
		cfg.Connections[name] = conn
		err = config.SaveConfig(cfg, GetPassphrase())
		if err != nil {
			fmt.Printf("\n⚠️  Key pushed but failed to update config: %v\n", err)
		} else {
			fmt.Printf("\n\033[32m✓\033[0m Public key successfully pushed and connection updated!\n")
		}

		fmt.Println("\033[90mYou can now connect to \033[1m" + name + "\033[0m\033[90m without a password.\033[0m\n")
	},
}

func init() {
	rootCmd.AddCommand(pushKeyCmd)
}
