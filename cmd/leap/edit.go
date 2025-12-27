package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit an existing SSH connection",
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
			fmt.Printf("\n❌ Connection '\033[1;36m%s\033[0m' not found.\n", name)
			fmt.Println("\033[90mTip: Use 'leap list' to see all available connections\033[0m\n")
			return
		}

		fmt.Println("\n⚡ \033[1;32mEdit SSH Connection\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\033[90mEditing: \033[1;36m%s\033[0m\n\n", name)

		promptHost := promptui.Prompt{
			Label:   "🌐 Hostname",
			Default: conn.Host,
		}
		host, _ := promptHost.Run()

		promptUser := promptui.Prompt{
			Label:   "👤 User",
			Default: conn.User,
		}
		user, _ := promptUser.Run()

		promptPort := promptui.Prompt{
			Label:   "🔌 Port",
			Default: fmt.Sprintf("%d", conn.Port),
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				return err
			},
		}
		portStr, _ := promptPort.Run()
		port, _ := strconv.Atoi(portStr)

		promptPass := promptui.Prompt{
			Label:   "🔐 Password (leave empty to keep current)",
			Mask:    '*',
			Default: "",
		}
		password, _ := promptPass.Run()
		if password == "" {
			password = conn.Password
		}

		promptKey := promptui.Prompt{
			Label:   "🔑 SSH Key Path",
			Default: conn.IdentityFile,
		}
		key, _ := promptKey.Run()

		promptTags := promptui.Prompt{
			Label:   "🏷️  Tags (comma separated)",
			Default: strings.Join(conn.Tags, ", "),
		}
		tagsStr, _ := promptTags.Run()
		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		promptJump := promptui.Prompt{
			Label:   "🔀 Jump Host",
			Default: conn.JumpHost,
		}
		jump, _ := promptJump.Run()

		promptNotes := promptui.Prompt{
			Label:   "📝 Notes",
			Default: conn.Notes,
		}
		notes, _ := promptNotes.Run()

		conn.Host = host
		conn.User = user
		conn.Port = port
		if password != "" {
			conn.Password = password
		}
		conn.IdentityFile = key
		conn.Tags = tags
		conn.JumpHost = jump
		conn.Notes = notes

		cfg.Connections[name] = conn

		err = config.SaveConfig(cfg, GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error saving config: %v\n\n", err)
			return
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Connection \033[1;36m%s\033[0m updated successfully!\n\n", name)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
