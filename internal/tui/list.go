package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paramientos/leap/internal/config"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			MarginLeft(2)

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#00BFFF")).
			Padding(0, 1).
			MarginRight(1).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#32CD32")).
			Padding(0, 1).
			Bold(true)
)

type item struct {
	title, desc string
	conn        config.Connection
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string {
	return i.title + " " + i.desc + " " + strings.Join(i.conn.Tags, " ")
}

type Model struct {
	list     list.Model
	choice   *config.Connection
	quitting bool
	width    int
	height   int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i.conn
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width/3-h, msg.Height-v-4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.choice != nil {
		return lipgloss.NewStyle().Padding(1, 2).Render(fmt.Sprintf("🚀 Connecting to %s...", m.choice.Name))
	}
	if m.quitting {
		return ""
	}

	curItem, ok := m.list.SelectedItem().(item)
	var details string
	if ok {
		conn := curItem.conn
		var tagsStr strings.Builder
		for _, tag := range conn.Tags {
			tagsStr.WriteString(tagStyle.Render(tag))
			tagsStr.WriteString(" ")
		}

		info := []string{
			fmt.Sprintf("Name:     %s", lipgloss.NewStyle().Bold(true).Render(conn.Name)),
			fmt.Sprintf("Host:     %s", conn.Host),
			fmt.Sprintf("User:     %s", conn.User),
			fmt.Sprintf("Port:     %d", conn.Port),
			fmt.Sprintf("Auth:     %s", getAuthType(conn)),
			"",
			fmt.Sprintf("Tags:     %s", tagsStr.String()),
			"",
			fmt.Sprintf("Status:   %s", statusStyle.Render("REACHABLE")),
		}

		dWidth := m.width - m.width/3 - 10
		if dWidth < 20 {
			dWidth = 20
		}
		dHeight := m.height - 8
		if dHeight < 5 {
			dHeight = 5
		}

		details = detailStyle.
			Width(dWidth).
			Height(dHeight).
			Render(strings.Join(info, "\n"))
	}

	listView := m.list.View()
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, details)

	return appStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("LEAP SSH MANAGER"),
			"",
			mainView,
		),
	)
}

func getAuthType(conn config.Connection) string {
	if conn.IdentityFile != "" {
		return "Key (" + conn.IdentityFile + ")"
	}
	if conn.Password != "" {
		return "Password (Saved)"
	}
	return "System Agent / Prompt"
}

func InitialModel(cfg *config.Config) Model {
	items := []list.Item{}
	for _, conn := range cfg.Connections {
		items = append(items, item{
			title: conn.Name,
			desc:  fmt.Sprintf("%s@%s", conn.User, conn.Host),
			conn:  conn,
		})
	}

	const defaultWidth = 40
	const listHeight = 14

	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "Connections"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)

	return Model{list: l}
}

func Run(cfg *config.Config) (*config.Connection, error) {
	m := InitialModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	res, ok := finalModel.(Model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}
	return res.choice, nil
}
