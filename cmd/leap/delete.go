package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [name...]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete SSH connection(s)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("\n❌ Please specify at least one connection name")
			fmt.Println("\033[90mUsage: leap delete [name1] [name2] ...\033[0m\n")
			return
		}

		cfg, err := config.LoadConfig(GetPassphrase())

		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		force, _ := cmd.Flags().GetBool("force")

		fmt.Println("\n⚡ \033[1;31mDelete SSH Connection(s)\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		deleted := []string{}
		notFound := []string{}

		for _, name := range args {
			if _, ok := cfg.Connections[name]; !ok {
				notFound = append(notFound, name)
				continue
			}

			if !force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Delete '%s'", name),
					IsConfirm: true,
				}
				result, err := prompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					fmt.Printf("\033[90m⊘ Skipped '%s'\033[0m\n", name)
					continue
				}
			}

			if cfg.DeleteConnection(name) {
				deleted = append(deleted, name)
				fmt.Printf("\033[32m✓\033[0m Deleted \033[1;36m%s\033[0m\n", name)
			}
		}

		if len(notFound) > 0 {
			fmt.Printf("\n\033[33m⚠\033[0m  Not found: %s\n", strings.Join(notFound, ", "))
		}

		if len(deleted) > 0 {
			err = config.SaveConfig(cfg, GetPassphrase())
			if err != nil {
				fmt.Printf("\n❌ Error saving config: %v\n\n", err)
				return
			}

			fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
			fmt.Printf("\n\033[32m✓\033[0m Successfully deleted \033[1m%d\033[0m connection(s)\n\n", len(deleted))
		} else {
			fmt.Println("\n\033[90mNo connections were deleted\033[0m\n")
		}
	},
}

func init() {
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	rootCmd.AddCommand(deleteCmd)
}
