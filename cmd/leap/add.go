package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new SSH connection",
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			prompt := promptui.Prompt{
				Label: "Connection Name (alias)",
			}
			var err error
			name, err = prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		}

		promptHost := promptui.Prompt{
			Label: "Hostname",
		}
		host, _ := promptHost.Run()

		promptUser := promptui.Prompt{
			Label:   "User",
			Default: "root",
		}
		user, _ := promptUser.Run()

		promptPort := promptui.Prompt{
			Label:   "Port",
			Default: "22",
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				return err
			},
		}
		portStr, _ := promptPort.Run()
		port, _ := strconv.Atoi(portStr)

		promptPass := promptui.Prompt{
			Label: "Password (optional)",
			Mask:  '*',
		}
		password, _ := promptPass.Run()

		promptKey := promptui.Prompt{
			Label: "SSH Key Path (optional)",
		}
		key, _ := promptKey.Run()

		promptTags := promptui.Prompt{
			Label: "Tags (comma separated)",
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
			Label: "Jump Host (optional)",
		}
		jump, _ := promptJump.Run()

		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		cfg.Connections[name] = config.Connection{
			Name:         name,
			Host:         host,
			User:         user,
			Port:         port,
			Password:     password,
			IdentityFile: key,
			Tags:         tags,
			JumpHost:     jump,
		}

		err = config.SaveConfig(cfg, GetPassphrase())
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("✓ Saved '%s'\n", name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
