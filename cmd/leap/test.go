package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test [name...]",
	Aliases: []string{"ping", "check"},
	Short:   "Test SSH connection(s)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(GetPassphrase())
		if err != nil {
			fmt.Printf("\n❌ Error loading config: %v\n\n", err)
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		tag, _ := cmd.Flags().GetString("tag")

		var connsToTest []config.Connection

		if all {
			for _, conn := range cfg.Connections {
				connsToTest = append(connsToTest, conn)
			}
		} else if tag != "" {
			for _, conn := range cfg.Connections {
				for _, t := range conn.Tags {
					if t == tag {
						connsToTest = append(connsToTest, conn)
						break
					}
				}
			}
		} else if len(args) > 0 {
			for _, name := range args {
				if conn, ok := cfg.Connections[name]; ok {
					connsToTest = append(connsToTest, conn)
				} else {
					fmt.Printf("\n\033[33m⚠\033[0m  Connection '\033[1;36m%s\033[0m' not found\n", name)
				}
			}
		} else {
			fmt.Println("\n❌ Please specify connection name(s), use --all, or --tag")
			fmt.Println("\033[90mUsage: leap test [name] or leap test --all\033[0m\n")
			return
		}

		if len(connsToTest) == 0 {
			fmt.Println("\n\033[90mNo connections to test\033[0m\n")
			return
		}

		fmt.Println("\n⚡ \033[1;32mConnection Health Check\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		for _, conn := range connsToTest {
			testConnection(conn)
		}

		fmt.Println("\n\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	},
}

func testConnection(conn config.Connection) {
	fmt.Printf("\033[1;36m%s\033[0m (\033[33m%s\033[0m@\033[32m%s\033[0m:\033[35m%d\033[0m)\n", conn.Name, conn.User, conn.Host, conn.Port)

	start := time.Now()
	address := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	tcpConn, err := net.DialTimeout("tcp", address, 5*time.Second)
	latency := time.Since(start)

	if err != nil {
		fmt.Printf("  \033[31m✗\033[0m Port %d: \033[31mUnreachable\033[0m (%v)\n", conn.Port, err)
		return
	}
	defer tcpConn.Close()

	fmt.Printf("  \033[32m✓\033[0m Port %d: \033[32mOpen\033[0m\n", conn.Port)
	fmt.Printf("  ⏱️  Latency: \033[36m%dms\033[0m\n", latency.Milliseconds())

	sshTestCmd := exec.Command("ssh",
		"-o", "ConnectTimeout=5",
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=no",
		"-p", fmt.Sprintf("%d", conn.Port),
		fmt.Sprintf("%s@%s", conn.User, conn.Host),
		"echo 'SSH OK'",
	)

	output, err := sshTestCmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "SSH OK") {
		fmt.Printf("  \033[32m✓\033[0m SSH Auth: \033[32mSuccess\033[0m\n")
	} else {
		fmt.Printf("  \033[33m⚠\033[0m  SSH Auth: \033[33mNot tested\033[0m (requires key/password)\n")
	}

	fmt.Println()
}

func init() {
	testCmd.Flags().BoolP("all", "a", false, "Test all connections")
	testCmd.Flags().StringP("tag", "t", "", "Test connections with specific tag")
	rootCmd.AddCommand(testCmd)
}
