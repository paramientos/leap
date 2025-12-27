package main

// This is the best module i think
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type ServerSnapshot struct {
	ServerName   string        `json:"server_name" yaml:"server_name"`
	Host         string        `json:"host" yaml:"host"`
	Timestamp    time.Time     `json:"timestamp" yaml:"timestamp"`
	OSInfo       OSInfo        `json:"os_info" yaml:"os_info"`
	SystemInfo   SystemInfo    `json:"system_info" yaml:"system_info"`
	Packages     []string      `json:"packages,omitempty" yaml:"packages,omitempty"`
	Services     []ServiceInfo `json:"services" yaml:"services"`
	OpenPorts    []string      `json:"open_ports" yaml:"open_ports"`
	DiskUsage    []DiskInfo    `json:"disk_usage" yaml:"disk_usage"`
	NetworkInfo  NetworkInfo   `json:"network_info" yaml:"network_info"`
	ProcessCount int           `json:"process_count" yaml:"process_count"`
	LoadAverage  string        `json:"load_average" yaml:"load_average"`
	Uptime       string        `json:"uptime" yaml:"uptime"`
}

type OSInfo struct {
	Distribution string `json:"distribution" yaml:"distribution"`
	Version      string `json:"version" yaml:"version"`
	Kernel       string `json:"kernel" yaml:"kernel"`
}

type SystemInfo struct {
	CPUCores     string `json:"cpu_cores" yaml:"cpu_cores"`
	TotalRAM     string `json:"total_ram" yaml:"total_ram"`
	UsedRAM      string `json:"used_ram" yaml:"used_ram"`
	Architecture string `json:"architecture" yaml:"architecture"`
}

type ServiceInfo struct {
	Name   string `json:"name" yaml:"name"`
	Status string `json:"status" yaml:"status"`
}

type DiskInfo struct {
	Filesystem string `json:"filesystem" yaml:"filesystem"`
	Size       string `json:"size" yaml:"size"`
	Used       string `json:"used" yaml:"used"`
	Available  string `json:"available" yaml:"available"`
	UsePercent string `json:"use_percent" yaml:"use_percent"`
	MountPoint string `json:"mount_point" yaml:"mount_point"`
}

type NetworkInfo struct {
	Interfaces []string `json:"interfaces" yaml:"interfaces"`
	PublicIP   string   `json:"public_ip" yaml:"public_ip"`
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot [connection]",
	Short: "Capture a comprehensive snapshot of a server's current state",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		output, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		includePackages, _ := cmd.Flags().GetBool("packages")

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

		fmt.Printf("\n📸 \033[1;32mCAPTURING SNAPSHOT\033[0m\n")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
		fmt.Printf("  Server: \033[1;36m%s\033[0m (%s@%s)\n\n", name, conn.User, conn.Host)

		snapshot, err := captureSnapshot(conn, name, includePackages)

		if err != nil {
			fmt.Printf("\n❌ Failed to capture snapshot: %v\n\n", err)
			return
		}

		var data []byte

		if format == "yaml" {
			data, err = yaml.Marshal(snapshot)
		} else {
			data, err = json.MarshalIndent(snapshot, "", "  ")
		}

		if err != nil {
			fmt.Printf("\n❌ Failed to marshal snapshot: %v\n\n", err)
			return
		}

		if output != "" {
			err = os.WriteFile(output, data, 0644)
			if err != nil {
				fmt.Printf("\n❌ Failed to write file: %v\n\n", err)
				return
			}
			fmt.Printf("\n\033[32m✓\033[0m Snapshot saved to \033[1;33m%s\033[0m\n\n", output)
		} else {
			fmt.Println(string(data))
		}
	},
}

func captureSnapshot(conn config.Connection, name string, includePackages bool) (*ServerSnapshot, error) {
	snapshot := &ServerSnapshot{
		ServerName: name,
		Host:       conn.Host,
		Timestamp:  time.Now(),
	}

	sshConfig := &ssh.ClientConfig{
		User:            conn.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	if conn.IdentityFile != "" {
		key, err := os.ReadFile(conn.IdentityFile)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
			}
		}
	}

	if conn.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(conn.Password))
	}

	addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)

	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	fmt.Print("  ⏳ Gathering OS information...")
	osInfo, _ := runCommand(client, "cat /etc/os-release | grep -E '^(NAME|VERSION)=' | head -2")
	kernelInfo, _ := runCommand(client, "uname -r")
	snapshot.OSInfo = parseOSInfo(osInfo, kernelInfo)
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Gathering system information...")
	cpuInfo, _ := runCommand(client, "nproc")
	ramInfo, _ := runCommand(client, "free -h | awk 'NR==2{print $2,$3}'")
	archInfo, _ := runCommand(client, "uname -m")
	snapshot.SystemInfo = parseSystemInfo(cpuInfo, ramInfo, archInfo)
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Gathering load and uptime...")
	loadAvg, _ := runCommand(client, "cat /proc/loadavg | awk '{print $1,$2,$3}'")
	uptime, _ := runCommand(client, "uptime -p")
	snapshot.LoadAverage = strings.TrimSpace(loadAvg)
	snapshot.Uptime = strings.TrimSpace(uptime)
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Gathering disk usage...")
	diskInfo, _ := runCommand(client, "df -h | tail -n +2")
	snapshot.DiskUsage = parseDiskInfo(diskInfo)
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Gathering service status...")
	services, _ := runCommand(client, "systemctl list-units --type=service --state=running --no-pager --no-legend | awk '{print $1}' | head -20")
	snapshot.Services = parseServices(services)
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Scanning open ports...")
	ports, _ := runCommand(client, "ss -tuln | grep LISTEN | awk '{print $5}' | sed 's/.*://' | sort -u")
	snapshot.OpenPorts = strings.Split(strings.TrimSpace(ports), "\n")
	fmt.Println(" \033[32m✓\033[0m")

	fmt.Print("  ⏳ Gathering network information...")
	interfaces, _ := runCommand(client, "ip -o link show | awk -F': ' '{print $2}' | grep -v lo")
	publicIP, _ := runCommand(client, "curl -s ifconfig.me || echo 'N/A'")
	snapshot.NetworkInfo = NetworkInfo{
		Interfaces: strings.Split(strings.TrimSpace(interfaces), "\n"),
		PublicIP:   strings.TrimSpace(publicIP),
	}
	fmt.Println(" \033[32m✓\033[0m")

	processCount, _ := runCommand(client, "ps aux | wc -l")
	fmt.Sscanf(processCount, "%d", &snapshot.ProcessCount)

	if includePackages {
		fmt.Print("  ⏳ Gathering installed packages (this may take a while)...")
		packages, _ := runCommand(client, "dpkg -l | grep ^ii | awk '{print $2}' || rpm -qa")
		snapshot.Packages = strings.Split(strings.TrimSpace(packages), "\n")
		fmt.Println(" \033[32m✓\033[0m")
	}

	return snapshot, nil
}

func runCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

func parseOSInfo(osRelease, kernel string) OSInfo {
	lines := strings.Split(osRelease, "\n")
	info := OSInfo{Kernel: strings.TrimSpace(kernel)}

	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			info.Distribution = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
		} else if strings.HasPrefix(line, "VERSION=") {
			info.Version = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
		}
	}
	return info
}

func parseSystemInfo(cpu, ram, arch string) SystemInfo {
	ramParts := strings.Fields(ram)
	info := SystemInfo{
		CPUCores:     strings.TrimSpace(cpu),
		Architecture: strings.TrimSpace(arch),
	}
	if len(ramParts) >= 2 {
		info.TotalRAM = ramParts[0]
		info.UsedRAM = ramParts[1]
	}
	return info
}

func parseDiskInfo(diskOutput string) []DiskInfo {
	var disks []DiskInfo
	lines := strings.Split(strings.TrimSpace(diskOutput), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			disks = append(disks, DiskInfo{
				Filesystem: fields[0],
				Size:       fields[1],
				Used:       fields[2],
				Available:  fields[3],
				UsePercent: fields[4],
				MountPoint: fields[5],
			})
		}
	}
	return disks
}

func parseServices(servicesOutput string) []ServiceInfo {
	var services []ServiceInfo
	lines := strings.Split(strings.TrimSpace(servicesOutput), "\n")

	for _, line := range lines {
		if line != "" {
			services = append(services, ServiceInfo{
				Name:   strings.TrimSpace(line),
				Status: "running",
			})
		}
	}
	return services
}

func init() {
	snapshotCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	snapshotCmd.Flags().StringP("format", "f", "json", "Output format: json or yaml")
	snapshotCmd.Flags().BoolP("packages", "p", false, "Include installed packages (slower)")

	rootCmd.AddCommand(snapshotCmd)
}
