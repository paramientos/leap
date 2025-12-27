package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "List and play recorded SSH sessions",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		historyDir := filepath.Join(home, ".leap", "history")

		if _, err := os.Stat(historyDir); os.IsNotExist(err) {
			fmt.Println("\n📜 No history found. Record a session with \033[1m'leap connect [name] --record'\033[0m")
			return
		}

		files, err := os.ReadDir(historyDir)
		if err != nil {
			fmt.Printf("\n❌ Error reading history: %v\n", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("\n📜 No recordings found.")
			return
		}

		fmt.Println("\n📜 \033[1;32mSESSION RECORDINGS\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")

		// Sort by date (filename contains date)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() > files[j].Name()
		})

		for i, f := range files {
			if strings.HasSuffix(f.Name(), ".cast") {
				info, _ := f.Info()
				dateStr := info.ModTime().Format("2006-01-02 15:04:05")
				name := strings.TrimSuffix(f.Name(), ".cast")

				fmt.Printf("  [%d] \033[1;36m%-20s\033[0m \033[90m%s\033[0m\n", i+1, name, dateStr)
			}
		}

		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Println("  Tip: Run \033[1m'leap replay [name_date]'\033[0m to view a recording.\n")
	},
}

var replayCmd = &cobra.Command{
	Use:   "replay [filename]",
	Short: "Replay a recorded session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, ".cast") {
			name += ".cast"
		}

		home, _ := os.UserHomeDir()
		path := filepath.Join(home, ".leap", "history", name)

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("\n❌ Error reading recording: %v\n\n", err)
			return
		}

		fmt.Printf("\n\033[1;33m▶ REPLAYING SESSION: %s\033[0m\n", name)
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		// For now, just print the content. In a future update, we can add real playback.
		fmt.Print(string(data))

		fmt.Println("\n\n\033[90m━━━━━━━━━━━━━━━━━━━━ END OF REPLAY ━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(replayCmd)
}
