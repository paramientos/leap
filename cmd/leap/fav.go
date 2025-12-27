package main

import (
	"fmt"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var favCmd = &cobra.Command{
	Use:     "favorite [name]",
	Aliases: []string{"fav"},
	Short:   "Toggle favorite status for a connection",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		isFav := cfg.ToggleFavorite(name)
		err = config.SaveConfig(cfg, GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error saving config: %v\n\n", err)
			return
		}

		if isFav {
			fmt.Printf("\n⭐ Connection \033[1;36m%s\033[0m marked as favorite!\n\n", name)
		} else {
			fmt.Printf("\n⚪ Connection \033[1;36m%s\033[0m removed from favorites.\n\n", name)
		}
	},
}

var listFavsCmd = &cobra.Command{
	Use:   "favorites",
	Short: "List all favorite connections",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		fmt.Println("\n⚡ \033[1;32mFavorite Connections\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		count := 0
		for name, conn := range cfg.Connections {
			if conn.Favorite {
				fmt.Printf(" ⭐ \033[1;36m%-20s\033[0m %s@%s:%d\n", name, conn.User, conn.Host, conn.Port)
				count++
			}
		}

		if count == 0 {
			fmt.Println("  No favorites yet. Use 'leap fav [name]' to add some!")
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("✓ Total favorites: %d\n\n", count)
	},
}

func init() {
	rootCmd.AddCommand(favCmd)
	rootCmd.AddCommand(listFavsCmd)
}
