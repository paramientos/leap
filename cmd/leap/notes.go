package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var notesCmd = &cobra.Command{
	Use:   "notes [name]",
	Short: "View or edit notes for a connection",
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

		edit, _ := cmd.Flags().GetBool("edit")

		if edit {
			prompt := promptui.Prompt{
				Label:   "📝 Notes",
				Default: conn.Notes,
			}
			notes, err := prompt.Run()
			if err != nil {
				fmt.Printf("\n❌ Prompt failed: %v\n\n", err)
				return
			}

			cfg.SetNotes(name, notes)
			err = config.SaveConfig(cfg, GetPassphrase())
			if err != nil {
				fmt.Printf("\n❌ Error saving config: %v\n\n", err)
				return
			}

			fmt.Printf("\n\033[32m✓\033[0m Notes updated for \033[1;36m%s\033[0m\n\n", name)
		} else {
			fmt.Println("\n⚡ \033[1;32mConnection Notes\033[0m")
			fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
			fmt.Printf("\n\033[1;36m%s\033[0m\n\n", name)

			if conn.Notes != "" {
				fmt.Printf("\033[90m%s\033[0m\n", conn.Notes)
			} else {
				fmt.Println("\033[90mNo notes available\033[0m")
				fmt.Println("\033[90mTip: Use 'leap notes [name] --edit' to add notes\033[0m")
			}

			fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
		}
	},
}

func init() {
	notesCmd.Flags().BoolP("edit", "e", false, "Edit notes")
	rootCmd.AddCommand(notesCmd)
}
