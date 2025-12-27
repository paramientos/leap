package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var importSSHCmd = &cobra.Command{
	Use:   "import-ssh",
	Short: "Import connections from ~/.ssh/config",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		sshConfigPath := filepath.Join(home, ".ssh", "config")

		if _, err := os.Stat(sshConfigPath); os.IsNotExist(err) {
			fmt.Printf("\n❌ SSH config not found at %s\n\n", sshConfigPath)
			return
		}

		file, err := os.Open(sshConfigPath)
		if err != nil {
			fmt.Printf("\n❌ Error opening SSH config: %v\n\n", err)
			return
		}
		defer file.Close()

		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading LEAP config: %v\n\n", err)
			return
		}

		fmt.Println("\n⚡ \033[1;32mImporting from SSH Config\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		scanner := bufio.NewScanner(file)
		var currentConn *config.Connection
		added := 0

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			key := strings.ToLower(parts[0])
			value := parts[1]

			if key == "host" {
				if currentConn != nil && currentConn.Name != "*" {
					cfg.Connections[currentConn.Name] = *currentConn
					fmt.Printf("\033[32m✓\033[0m Added \033[1;36m%s\033[0m (%s)\n", currentConn.Name, currentConn.Host)
					added++
				}
				currentConn = &config.Connection{
					Name: value,
					Port: 22,
					User: "root",
				}
			} else if currentConn != nil {
				switch key {
				case "hostname":
					currentConn.Host = value
				case "user":
					currentConn.User = value
				case "port":
					port, _ := strconv.Atoi(value)
					currentConn.Port = port
				case "identityfile":
					currentConn.IdentityFile = value
				case "proxyjump":
					currentConn.JumpHost = value
				}
			}
		}

		// Final one
		if currentConn != nil && currentConn.Name != "*" {
			cfg.Connections[currentConn.Name] = *currentConn
			fmt.Printf("\033[32m✓\033[0m Added \033[1;36m%s\033[0m (%s)\n", currentConn.Name, currentConn.Host)
			added++
		}

		if added > 0 {
			err = config.SaveConfig(cfg, GetPassphrase())
			if err != nil {
				fmt.Printf("\n❌ Error saving config: %v\n\n", err)
				return
			}
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Successfully imported \033[1m%d\033[0m connections!\n\n", added)
	},
}

func init() {
	rootCmd.AddCommand(importSSHCmd)
}
