package main

import (
	"fmt"
	"os"

	"github.com/paramientos/leap/internal/config"
	"github.com/paramientos/leap/internal/tui"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor [name...]",
	Aliases: []string{"watch", "top"},
	Short:   "Monitor server resources (CPU, RAM, Load) in real-time",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		tag, _ := cmd.Flags().GetString("tag")

		var connsToMonitor []config.Connection

		if all {
			for _, conn := range cfg.Connections {
				connsToMonitor = append(connsToMonitor, conn)
			}
		} else if tag != "" {
			for _, conn := range cfg.Connections {
				for _, t := range conn.Tags {
					if t == tag {
						connsToMonitor = append(connsToMonitor, conn)
						break
					}
				}
			}
		} else if len(args) > 0 {
			for _, name := range args {
				if conn, ok := cfg.Connections[name]; ok {
					connsToMonitor = append(connsToMonitor, conn)
				} else {
					fmt.Printf("\n\033[33m⚠\033[0m  Connection '\033[1;36m%s\033[0m' not found\n", name)
				}
			}
		} else {
			// If no args, maybe just monitor everything or show selection?
			// Let's monitor everything by default if it's a TUI tool
			for _, conn := range cfg.Connections {
				connsToMonitor = append(connsToMonitor, conn)
			}
		}

		if len(connsToMonitor) == 0 {
			fmt.Println("\n\033[90mNo connections to monitor\033[0m\n")
			return
		}

		err = tui.RunMonitor(connsToMonitor)
		if err != nil {
			fmt.Printf("\n❌ Error running monitor: %v\n\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	monitorCmd.Flags().BoolP("all", "a", false, "Monitor all connections")
	monitorCmd.Flags().StringP("tag", "t", "", "Monitor connections with specific tag")
	rootCmd.AddCommand(monitorCmd)
}
