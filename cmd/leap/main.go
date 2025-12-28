package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/paramientos/leap/internal/config"
	"github.com/paramientos/leap/internal/ssh"
	"github.com/paramientos/leap/internal/tui"
	"github.com/paramientos/leap/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leap",
	Short: "⚡ LEAP - Modern SSH Connection Manager",
	Long: `
⚡ LEAP SSH MANAGER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

A modern CLI tool to manage your SSH connections with tags, 
fuzzy search, and an intuitive terminal interface.

Features:
  • 🔐 Secure encrypted configuration
  • 🏷️  Tag-based organization
  • 🔍 Fuzzy search & filtering
  • 🎨 Beautiful terminal UI
  • 🔀 Jump host support
  • 🚇 SSH tunnel management
`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		if len(args) > 0 {
			name := strings.Join(args, " ")
			if conn, ok := cfg.Connections[name]; ok {
				fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", name)
				err := ssh.Connect(conn, false)
				if err != nil {
					fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
				}
				return
			}

			for _, conn := range cfg.Connections {
				if strings.Contains(strings.ToLower(conn.Name), strings.ToLower(name)) {
					fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", conn.Name)
					err := ssh.Connect(conn, false)
					if err != nil {
						fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
					}
					return
				}

				for _, tag := range conn.Tags {
					if strings.EqualFold(tag, name) {
						fmt.Printf("\n🚀 Connecting to \033[1;36m%s\033[0m...\n\n", conn.Name)
						err := ssh.Connect(conn, false)
						if err != nil {
							fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
						}
						return
					}
				}
			}
		}

		choice, err := tui.Run(cfg)

		if err != nil {
			fmt.Printf("\n❌ Error running TUI: %v\n\n", err)
			return
		}

		if choice != nil {
			err = ssh.Connect(*choice, false)
			if err != nil {
				fmt.Printf("\n❌ SSH Connection closed with error: %v\n\n", err)
			}
		}
	},
}

var masterPassword string

func GetPassphrase() string {
	if masterPassword != "" {
		return masterPassword
	}

	envPass := os.Getenv("LEAP_MASTER_PASSWORD")
	if envPass != "" {
		masterPassword = envPass
		return masterPassword
	}

	home, _ := os.UserHomeDir()
	sessionFile := filepath.Join(home, ".leap", ".session")
	hostname, _ := os.Hostname()
	salt := hostname + home // Unique to this machine and user

	// Check for active session (5 minute cache)
	if info, err := os.Stat(sessionFile); err == nil {
		if time.Since(info.ModTime()) < 5*time.Minute {
			data, err := os.ReadFile(sessionFile)
			if err == nil {
				decrypted, err := utils.Deobfuscate(data, salt)
				if err == nil {
					masterPassword = string(decrypted)
					// Refresh the session timeout on each use
					os.Chtimes(sessionFile, time.Now(), time.Now())
					return masterPassword
				} else {
					// File exists but is not valid obfuscated data (maybe old plaintext)
					// Delete it to be safe
					os.Remove(sessionFile)
				}
			}
		} else {
			// Session expired
			os.Remove(sessionFile)
		}
	}

	path := config.GetConfigPath()
	isFirstRun := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isFirstRun = true
	}

	if isFirstRun {
		fmt.Println("\n✨ \033[1;32mWelcome to LEAP SSH Manager!\033[0m")
		fmt.Println("This is your first run. Let's set up a \033[1mMaster Password\033[0m to encrypt your connections.")
		fmt.Println("\033[90mNote: You can avoid this prompt by setting LEAP_MASTER_PASSWORD environment variable.\033[0m")

		prompt := promptui.Prompt{
			Label: "🔒 Set Master Password",
			Mask:  '*',
			Validate: func(input string) error {
				if len(input) < 4 {
					return fmt.Errorf("password must be at least 4 characters")
				}
				return nil
			},
		}
		res, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}
		masterPassword = res

		cfg := &config.Config{Connections: make(map[string]config.Connection)}
		err = config.SaveConfig(cfg, masterPassword)
		if err != nil {
			fmt.Printf("❌ Failed to initialize config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\033[32m✓ Master Password set and config initialized.\033[0m")
	} else {
		prompt := promptui.Prompt{
			Label: "🔓 Enter Master Password",
			Mask:  '*',
		}
		res, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}
		masterPassword = res
	}

	// Save to session cache (600 permissions - user only)
	encrypted, err := utils.Obfuscate([]byte(masterPassword), salt)
	if err == nil {
		os.WriteFile(sessionFile, encrypted, 0600)
	}

	return masterPassword
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("⚡ LEAP SSH Manager v{{.Version}}\n")
}
