package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [name] [local-path] [remote-path]",
	Short: "Upload file(s) to remote server",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		localPath := args[1]
		remotePath := args[2]

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

		fmt.Println("\n⚡ \033[1;32mFile Upload\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[90mFrom:\033[0m \033[1;35m%s\033[0m\n", localPath)
		fmt.Printf("\033[90mTo:\033[0m   \033[1;36m%s\033[0m:\033[1;35m%s\033[0m\n\n", conn.Name, remotePath)

		scpArgs := []string{}

		if conn.IdentityFile != "" {
			scpArgs = append(scpArgs, "-i", conn.IdentityFile)
		}

		scpArgs = append(scpArgs, "-P", fmt.Sprintf("%d", conn.Port))

		if conn.JumpHost != "" {
			scpArgs = append(scpArgs, "-o", fmt.Sprintf("ProxyJump=%s", conn.JumpHost))
		}

		recursive, _ := cmd.Flags().GetBool("recursive")
		if recursive {
			scpArgs = append(scpArgs, "-r")
		}

		scpArgs = append(scpArgs, localPath)
		scpArgs = append(scpArgs, fmt.Sprintf("%s@%s:%s", conn.User, conn.Host, remotePath))

		scpCmd := exec.Command("scp", scpArgs...)
		scpCmd.Stdout = os.Stdout
		scpCmd.Stderr = os.Stderr
		scpCmd.Stdin = os.Stdin

		err = scpCmd.Run()
		if err != nil {
			fmt.Printf("\n❌ Upload failed: %v\n\n", err)
			return
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Upload completed successfully\n\n")
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [name] [remote-path] [local-path]",
	Short: "Download file(s) from remote server",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		remotePath := args[1]
		localPath := args[2]

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

		fmt.Println("\n⚡ \033[1;32mFile Download\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[90mFrom:\033[0m \033[1;36m%s\033[0m:\033[1;35m%s\033[0m\n", conn.Name, remotePath)
		fmt.Printf("\033[90mTo:\033[0m   \033[1;35m%s\033[0m\n\n", localPath)

		scpArgs := []string{}

		if conn.IdentityFile != "" {
			scpArgs = append(scpArgs, "-i", conn.IdentityFile)
		}

		scpArgs = append(scpArgs, "-P", fmt.Sprintf("%d", conn.Port))

		if conn.JumpHost != "" {
			scpArgs = append(scpArgs, "-o", fmt.Sprintf("ProxyJump=%s", conn.JumpHost))
		}

		recursive, _ := cmd.Flags().GetBool("recursive")
		if recursive {
			scpArgs = append(scpArgs, "-r")
		}

		scpArgs = append(scpArgs, fmt.Sprintf("%s@%s:%s", conn.User, conn.Host, remotePath))
		scpArgs = append(scpArgs, localPath)

		scpCmd := exec.Command("scp", scpArgs...)
		scpCmd.Stdout = os.Stdout
		scpCmd.Stderr = os.Stderr
		scpCmd.Stdin = os.Stdin

		err = scpCmd.Run()
		if err != nil {
			fmt.Printf("\n❌ Download failed: %v\n\n", err)
			return
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Download completed successfully\n\n")
	},
}

func init() {
	uploadCmd.Flags().BoolP("recursive", "r", false, "Upload directory recursively")
	downloadCmd.Flags().BoolP("recursive", "r", false, "Download directory recursively")

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
}
