package tui

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type stats struct {
	CPU    string
	RAM    string
	Load   string
	Uptime string
	Error  error
}

type serverStats struct {
	index int
	stats stats
}

type tickMsg time.Time

type monitorModel struct {
	connections []config.Connection
	stats       map[int]stats
	table       table.Model
	quitting    bool
	width       int
	height      int
}

func (m monitorModel) Init() tea.Cmd {
	return tea.Batch(
		m.updateAllStats(),
		tick(),
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m monitorModel) updateAllStats() tea.Cmd {
	var cmds []tea.Cmd
	for i, conn := range m.connections {
		cmds = append(cmds, m.fetchStats(i, conn))
	}
	return tea.Batch(cmds...)
}

func (m monitorModel) fetchStats(index int, conn config.Connection) tea.Cmd {
	return func() tea.Msg {
		s := stats{}

		sshConfig := &ssh.ClientConfig{
			User:            conn.User,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		// Try ID File
		if conn.IdentityFile != "" {
			key, err := os.ReadFile(conn.IdentityFile)
			if err == nil {
				signer, err := ssh.ParsePrivateKey(key)
				if err == nil {
					sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
				}
			}
		}

		// Support SSH Agent
		if socket := os.Getenv("SSH_AUTH_SOCK"); socket != "" {
			netConn, err := net.Dial("unix", socket)
			if err == nil {
				sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agent.NewClient(netConn).Signers))
			}
		}

		// Try saved password
		if conn.Password != "" {
			sshConfig.Auth = append(sshConfig.Auth, ssh.Password(conn.Password))
		}

		addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
		client, err := ssh.Dial("tcp", addr, sshConfig)
		if err != nil {
			s.Error = err
			return serverStats{index, s}
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			s.Error = err
			return serverStats{index, s}
		}
		defer session.Close()

		// Robust command: each output on its own line
		cmd := "cat /proc/uptime | awk '{print int($1/3600)\"h \"int(($1%3600)/60)\"m\"}'; " +
			"free | awk 'NR==2{print int($3*100/$2)\"%\"}'; " +
			"cat /proc/loadavg | awk '{print $1}'"

		output, err := session.Output(cmd)
		if err != nil {
			s.Error = err
			return serverStats{index, s}
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) >= 3 {
			s.Uptime = lines[0]
			s.RAM = lines[1]
			s.Load = lines[2]
		} else {
			s.Error = fmt.Errorf("unexpected output format")
		}

		return serverStats{index, s}
	}
}

func (m monitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Batch(m.updateAllStats(), tick())
	case serverStats:
		m.stats[msg.index] = msg.stats
		m.updateTableRows()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(msg.Width - 4)
		// Leave some room for header/footer
		m.table.SetHeight(msg.Height - 10)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *monitorModel) updateTableRows() {
	rows := []table.Row{}
	for i, conn := range m.connections {
		s, ok := m.stats[i]

		name := conn.Name
		if conn.Favorite {
			name = "⭐ " + name
		}

		if !ok {
			rows = append(rows, table.Row{name, "...", "...", "...", "CONNECTING..."})
			continue
		}

		if s.Error != nil {
			errStr := s.Error.Error()
			if len(errStr) > 20 {
				errStr = errStr[:17] + "..."
			}
			rows = append(rows, table.Row{name, "ERR", "ERR", "ERR", "❌ " + errStr})
		} else {
			rows = append(rows, table.Row{
				name,
				s.Load,
				s.RAM,
				s.Uptime,
				"✅ ONLINE",
			})
		}
	}
	m.table.SetRows(rows)
}

func (m monitorModel) View() string {
	if m.quitting {
		return ""
	}

	header := headerStyle.Render("📊 LIVE SERVER MONITOR")
	subtitle := subtitleStyle.Render(fmt.Sprintf("Real-time stats for %d connections • Updates every 5s", len(m.connections)))

	// Create a nice box for the table
	tableBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0).
		Render(m.table.View())

	footer := helpStyle.Render("Press q to exit • Select a row and press Enter to connect 🚀")

	return appStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			subtitle,
			"",
			tableBox,
			"",
			footer,
		),
	)
}

func RunMonitor(conns []config.Connection) error {
	columns := []table.Column{
		{Title: "SERVER", Width: 18},
		{Title: "LOAD", Width: 10},
		{Title: "RAM", Width: 10},
		{Title: "UPTIME", Width: 15},
		{Title: "STATUS", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Foreground(lipgloss.Color("86")).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("62")).
		Bold(true)
	t.SetStyles(s)

	m := monitorModel{
		connections: conns,
		stats:       make(map[int]stats),
		table:       t,
	}

	m.updateTableRows()

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
