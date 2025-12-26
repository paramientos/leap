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
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		tagFilter, _ := cmd.Flags().GetString("tag")

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "NAME\tCONNECTION\tTAGS")

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

			connectionStr := fmt.Sprintf("%s@%s:%d", conn.User, conn.Host, conn.Port)
			tagsStr := "[" + strings.Join(conn.Tags, ",") + "]"
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, connectionStr, tagsStr)
		}
		w.Flush()
	},
}

func init() {
	listCmd.Flags().StringP("tag", "t", "", "Filter by tag")
	rootCmd.AddCommand(listCmd)
}
