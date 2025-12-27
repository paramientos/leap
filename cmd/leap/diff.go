package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [snapshot1] [snapshot2]",
	Short: "Compare two server snapshots and show differences",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		file2 := args[1]

		fmt.Println("\n🔍 \033[1;32mCOMPARING SNAPSHOTS\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		snap1, err := loadSnapshot(file1)
		if err != nil {
			fmt.Printf("❌ Error loading %s: %v\n\n", file1, err)
			return
		}

		snap2, err := loadSnapshot(file2)
		if err != nil {
			fmt.Printf("❌ Error loading %s: %v\n\n", file2, err)
			return
		}

		fmt.Printf("  📄 Snapshot 1: \033[1;36m%s\033[0m (%s)\n", snap1.ServerName, snap1.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  📄 Snapshot 2: \033[1;36m%s\033[0m (%s)\n\n", snap2.ServerName, snap2.Timestamp.Format("2006-01-02 15:04:05"))

		hasChanges := false

		// Compare OS Info
		if snap1.OSInfo != snap2.OSInfo {
			hasChanges = true
			fmt.Println("  \033[1;33m▸ OS Information Changed\033[0m")
			if snap1.OSInfo.Kernel != snap2.OSInfo.Kernel {
				fmt.Printf("    Kernel: \033[31m%s\033[0m → \033[32m%s\033[0m\n", snap1.OSInfo.Kernel, snap2.OSInfo.Kernel)
			}
			if snap1.OSInfo.Version != snap2.OSInfo.Version {
				fmt.Printf("    Version: \033[31m%s\033[0m → \033[32m%s\033[0m\n", snap1.OSInfo.Version, snap2.OSInfo.Version)
			}
			fmt.Println()
		}

		// Compare System Resources
		if snap1.SystemInfo != snap2.SystemInfo {
			hasChanges = true
			fmt.Println("  \033[1;33m▸ System Resources Changed\033[0m")
			if snap1.SystemInfo.UsedRAM != snap2.SystemInfo.UsedRAM {
				fmt.Printf("    RAM Usage: \033[31m%s\033[0m → \033[32m%s\033[0m\n", snap1.SystemInfo.UsedRAM, snap2.SystemInfo.UsedRAM)
			}
			fmt.Println()
		}

		// Compare Load Average
		if snap1.LoadAverage != snap2.LoadAverage {
			hasChanges = true
			fmt.Println("  \033[1;33m▸ Load Average Changed\033[0m")
			fmt.Printf("    \033[31m%s\033[0m → \033[32m%s\033[0m\n\n", snap1.LoadAverage, snap2.LoadAverage)
		}

		// Compare Services
		services1 := make(map[string]bool)
		services2 := make(map[string]bool)

		for _, s := range snap1.Services {
			services1[s.Name] = true
		}
		for _, s := range snap2.Services {
			services2[s.Name] = true
		}

		var newServices, removedServices []string
		for name := range services2 {
			if !services1[name] {
				newServices = append(newServices, name)
			}
		}
		for name := range services1 {
			if !services2[name] {
				removedServices = append(removedServices, name)
			}
		}

		if len(newServices) > 0 || len(removedServices) > 0 {
			hasChanges = true
			fmt.Println("  \033[1;33m▸ Services Changed\033[0m")
			if len(newServices) > 0 {
				fmt.Printf("    \033[32m+ New:\033[0m %s\n", strings.Join(newServices, ", "))
			}
			if len(removedServices) > 0 {
				fmt.Printf("    \033[31m- Removed:\033[0m %s\n", strings.Join(removedServices, ", "))
			}
			fmt.Println()
		}

		// Compare Open Ports
		ports1 := make(map[string]bool)
		ports2 := make(map[string]bool)

		for _, p := range snap1.OpenPorts {
			ports1[p] = true
		}
		for _, p := range snap2.OpenPorts {
			ports2[p] = true
		}

		var newPorts, closedPorts []string
		for port := range ports2 {
			if !ports1[port] && port != "" {
				newPorts = append(newPorts, port)
			}
		}
		for port := range ports1 {
			if !ports2[port] && port != "" {
				closedPorts = append(closedPorts, port)
			}
		}

		if len(newPorts) > 0 || len(closedPorts) > 0 {
			hasChanges = true
			fmt.Println("  \033[1;33m▸ Open Ports Changed\033[0m")
			if len(newPorts) > 0 {
				fmt.Printf("    \033[32m+ Opened:\033[0m %s\n", strings.Join(newPorts, ", "))
			}
			if len(closedPorts) > 0 {
				fmt.Printf("    \033[31m- Closed:\033[0m %s\n", strings.Join(closedPorts, ", "))
			}
			fmt.Println()
		}

		// Compare Disk Usage
		if len(snap1.DiskUsage) > 0 && len(snap2.DiskUsage) > 0 {
			disk1 := snap1.DiskUsage[0]
			disk2 := snap2.DiskUsage[0]

			if disk1.UsePercent != disk2.UsePercent {
				hasChanges = true
				fmt.Println("  \033[1;33m▸ Disk Usage Changed\033[0m")
				fmt.Printf("    %s: \033[31m%s\033[0m → \033[32m%s\033[0m\n\n", disk1.MountPoint, disk1.UsePercent, disk2.UsePercent)
			}
		}

		// Compare Packages (if available)
		if len(snap1.Packages) > 0 && len(snap2.Packages) > 0 {
			packages1 := make(map[string]bool)
			packages2 := make(map[string]bool)

			for _, p := range snap1.Packages {
				packages1[p] = true
			}
			for _, p := range snap2.Packages {
				packages2[p] = true
			}

			var newPackages, removedPackages []string
			for pkg := range packages2 {
				if !packages1[pkg] && pkg != "" {
					newPackages = append(newPackages, pkg)
				}
			}
			for pkg := range packages1 {
				if !packages2[pkg] && pkg != "" {
					removedPackages = append(removedPackages, pkg)
				}
			}

			if len(newPackages) > 0 || len(removedPackages) > 0 {
				hasChanges = true
				fmt.Println("  \033[1;33m▸ Packages Changed\033[0m")
				if len(newPackages) > 0 {
					fmt.Printf("    \033[32m+ Installed:\033[0m %d packages\n", len(newPackages))
					if len(newPackages) <= 10 {
						fmt.Printf("      %s\n", strings.Join(newPackages, ", "))
					}
				}
				if len(removedPackages) > 0 {
					fmt.Printf("    \033[31m- Removed:\033[0m %d packages\n", len(removedPackages))
					if len(removedPackages) <= 10 {
						fmt.Printf("      %s\n", strings.Join(removedPackages, ", "))
					}
				}
				fmt.Println()
			}
		}

		if !hasChanges {
			fmt.Println("  \033[32m✓ No significant changes detected\033[0m\n")
		}

		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	},
}

func loadSnapshot(filename string) (*ServerSnapshot, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var snapshot ServerSnapshot
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
