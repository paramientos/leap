package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SSH connections",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			return
		}

		tagFilter, _ := cmd.Flags().GetString("tag")

		// Print header
		fmt.Println("\n⚡ \033[1;32mLEAP SSH MANAGER\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")

		if tagFilter != "" {
			fmt.Printf("\033[36mFiltered by tag:\033[0m \033[1;33m#%s\033[0m\n\n", tagFilter)
		} else {
			fmt.Println()
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 3, ' ', 0)
		fmt.Fprintln(w, "\033[1;36mNAME\033[0m\t\033[1;36mCONNECTION\033[0m\t\033[1;36mTAGS\033[0m")
		fmt.Fprintln(w, "\033[90m────\033[0m\t\033[90m──────────\033[0m\t\033[90m────\033[0m")

		count := 0
		for name, conn := range cfg.Connections {
			if tagFilter != "" {
				found := false
				for _, t := range conn.Tags {
					if t == tagFilter {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			connectionStr := fmt.Sprintf("\033[33m%s\033[0m@\033[32m%s\033[0m:\033[35m%d\033[0m", conn.User, conn.Host, conn.Port)

			var tagsStr string
			if len(conn.Tags) > 0 {
				tagParts := []string{}
				for _, tag := range conn.Tags {
					tagParts = append(tagParts, "\033[36m#"+tag+"\033[0m")
				}
				tagsStr = strings.Join(tagParts, " ")
			} else {
				tagsStr = "\033[90m-\033[0m"
			}

			fmt.Fprintf(w, "\033[1m%s\033[0m\t%s\t%s\n", name, connectionStr, tagsStr)
			count++
		}
		w.Flush()

		fmt.Println()
		fmt.Printf("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
		fmt.Printf("\033[32m✓\033[0m Total connections: \033[1m%d\033[0m\n\n", count)
	},
}

func init() {
	listCmd.Flags().StringP("tag", "t", "", "Filter by tag")
	rootCmd.AddCommand(listCmd)
}
