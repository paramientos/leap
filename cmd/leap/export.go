package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var exportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export connections to a file",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		format, _ := cmd.Flags().GetString("format")
		var data []byte

		switch format {
		case "json":
			data, err = json.MarshalIndent(cfg, "", "  ")
		case "yaml":
			data, err = yaml.Marshal(cfg)
		default:
			fmt.Printf("\n❌ Unknown format: %s (use 'json' or 'yaml')\n\n", format)
			return
		}

		if err != nil {
			fmt.Printf("\n❌ Error marshaling config: %v\n\n", err)
			return
		}

		if len(args) > 0 {
			filename := args[0]
			err = os.WriteFile(filename, data, 0600)
			if err != nil {
				fmt.Printf("\n❌ Error writing file: %v\n\n", err)
				return
			}
			fmt.Printf("\n\033[32m✓\033[0m Exported to \033[1;36m%s\033[0m (%s format)\n\n", filename, format)
		} else {
			fmt.Println(string(data))
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import [filename]",
	Short: "Import connections from a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("\n❌ Error reading file: %v\n\n", err)
			return
		}

		var importedCfg config.Config
		format, _ := cmd.Flags().GetString("format")

		switch format {
		case "json":
			err = json.Unmarshal(data, &importedCfg)
		case "yaml":
			err = yaml.Unmarshal(data, &importedCfg)
		default:
			err = json.Unmarshal(data, &importedCfg)
			if err != nil {
				err = yaml.Unmarshal(data, &importedCfg)
			}
		}

		if err != nil {
			fmt.Printf("\n❌ Error parsing file: %v\n\n", err)
			return
		}

		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		merge, _ := cmd.Flags().GetBool("merge")
		added := 0
		updated := 0
		skipped := 0

		fmt.Println("\n⚡ \033[1;32mImport Connections\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		for name, conn := range importedCfg.Connections {
			if _, exists := cfg.Connections[name]; exists {
				if merge {
					cfg.Connections[name] = conn
					fmt.Printf("\033[33m⟳\033[0m Updated \033[1;36m%s\033[0m\n", name)
					updated++
				} else {
					fmt.Printf("\033[90m⊘ Skipped \033[1;36m%s\033[0m (already exists)\n", name)
					skipped++
				}
			} else {
				cfg.Connections[name] = conn
				fmt.Printf("\033[32m✓\033[0m Added \033[1;36m%s\033[0m\n", name)
				added++
			}
		}

		if added > 0 || updated > 0 {
			err = config.SaveConfig(cfg, GetPassphrase())
			if err != nil {
				fmt.Printf("\n❌ Error saving config: %v\n\n", err)
				return
			}
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Printf("\n\033[32m✓\033[0m Added: \033[1m%d\033[0m  ", added)
		fmt.Printf("\033[33m⟳\033[0m Updated: \033[1m%d\033[0m  ", updated)
		fmt.Printf("\033[90m⊘ Skipped: \033[1m%d\033[0m\n\n", skipped)
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "json", "Export format (json or yaml)")
	importCmd.Flags().StringP("format", "f", "auto", "Import format (json, yaml, or auto)")
	importCmd.Flags().BoolP("merge", "m", false, "Merge and update existing connections")

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}
