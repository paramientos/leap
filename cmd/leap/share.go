package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [name]",
	Short: "Share a connection via QR Code or Short-link",
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
			fmt.Printf("\n❌ Connection \033[1;36m%s\033[0m not found.\n\n", name)
			return
		}

		// Share only essential fields to keep QR code small
		type ShareData struct {
			N string   `json:"n"`
			H string   `json:"h"`
			U string   `json:"u"`
			P int      `json:"p"`
			T []string `json:"t,omitempty"`
			J string   `json:"j,omitempty"`
			G string   `json:"g,omitempty"`
		}

		shared := ShareData{
			N: conn.Name,
			H: conn.Host,
			U: conn.User,
			P: conn.Port,
			T: conn.Tags,
			J: conn.JumpHost,
			G: conn.Group,
		}

		jsonData, _ := json.Marshal(shared)
		encoded := base64.StdEncoding.EncodeToString(jsonData)

		fmt.Println("\n⚡ \033[1;32mSHARE CONNECTION\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		fmt.Printf("  Sharing: \033[1;36m%s\033[0m (%s@%s)\n", conn.Name, conn.User, conn.Host)
		fmt.Println("\n  \033[1mOption 1: Scan QR Code\033[0m")

		qrConfig := qrterminal.Config{
			Level:      qrterminal.L,
			Writer:     os.Stdout,
			HalfBlocks: true,
			QuietZone:  1,
		}
		qrterminal.GenerateWithConfig(encoded, qrConfig)

		fmt.Println("\n  \033[1mOption 2: Copy-Paste Short Code\033[0m")
		fmt.Printf("\n\033[90m  %s\033[0m\n", encoded)

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Println("  Tip: Receiver can run \033[1m'leap import-code [code]'\033[0m to add it.")
		fmt.Println("\n\033[90mNote: Passwords and local identity paths are removed for security.\033[0m\n")
	},
}

var importCodeCmd = &cobra.Command{
	Use:   "import-code [code]",
	Short: "Import a shared connection using a base64 code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := args[0]

		decoded, err := base64.StdEncoding.DecodeString(code)
		if err != nil {
			fmt.Printf("\n❌ Invalid share code: %v\n\n", err)
			return
		}

		// Re-define internal struct for decoding
		type shareData struct {
			N string   `json:"n"`
			H string   `json:"h"`
			U string   `json:"u"`
			P int      `json:"p"`
			T []string `json:"t"`
			J string   `json:"j"`
			G string   `json:"g"`
		}

		var s shareData
		err = json.Unmarshal(decoded, &s)
		if err != nil {
			fmt.Printf("\n❌ Error parsing connection data: %v\n\n", err)
			return
		}

		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		// Map back to Connection
		conn := config.Connection{
			Name:     s.N,
			Host:     s.H,
			User:     s.U,
			Port:     s.P,
			Tags:     s.T,
			JumpHost: s.J,
			Group:    s.G,
		}

		cfg.Connections[conn.Name] = conn
		err = config.SaveConfig(cfg, GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error saving config: %v\n\n", err)
			return
		}

		fmt.Printf("\n\033[32m✓\033[0m Connection \033[1;36m%s\033[0m imported successfully!\n", conn.Name)
		fmt.Println("\033[90mTip: Use 'leap push-key " + conn.Name + "' to setup passwordless login.\033[0m\n")
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(importCodeCmd)
}
