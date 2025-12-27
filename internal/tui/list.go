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
	laravelRed     = lipgloss.Color("#FF2D20")
	laravelOrange  = lipgloss.Color("#FF6B35")
	primaryGreen   = lipgloss.Color("#10B981")
	accentCyan     = lipgloss.Color("#06B6D4")
	darkBg         = lipgloss.Color("#1F2937")
	lightText      = lipgloss.Color("#F9FAFB")
	mutedText      = lipgloss.Color("#9CA3AF")
	borderColor    = lipgloss.Color("#374151")
	highlightColor = lipgloss.Color("#3B82F6")

	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(lipgloss.Color("#111827"))

	// Header with gradient effect (simulated)
	headerStyle = lipgloss.NewStyle().
			Foreground(lightText).
			Background(primaryGreen).
			Padding(0, 2).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedText).
			Italic(true).
			MarginBottom(1)
	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(2, 3).
			MarginLeft(2).
			Background(darkBg)

	labelStyle = lipgloss.NewStyle().
			Foreground(accentCyan).
			Bold(true).
			Width(12)
	valueStyle = lipgloss.NewStyle().
			Foreground(lightText)

	tagStyle = lipgloss.NewStyle().
			Foreground(darkBg).
			Background(accentCyan).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(darkBg).
			Background(primaryGreen).
			Padding(0, 2).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(darkBg).
			Background(laravelOrange).
			Padding(0, 2).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedText).
			MarginTop(1)

	connectionStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
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
		connectMsg := lipgloss.NewStyle().
			Foreground(lightText).
			Bold(true).
			Render(fmt.Sprintf("🚀 Connecting to %s...", connectionStyle.Render(m.choice.Name)))
		return appStyle.Render(connectMsg)
	}

	if m.quitting {
		return ""
	}

	curItem, ok := m.list.SelectedItem().(item)
	var details string

	if ok {
		conn := curItem.conn

		var tagsStr strings.Builder

		if len(conn.Tags) > 0 {
			for _, tag := range conn.Tags {
				tagsStr.WriteString(tagStyle.Render("# " + tag))
				tagsStr.WriteString(" ")
			}
		} else {
			tagsStr.WriteString(lipgloss.NewStyle().Foreground(mutedText).Render("No tags"))
		}

		authType := getAuthType(conn)
		authIcon := "🔑"

		if strings.Contains(authType, "Password") {
			authIcon = "🔐"
		} else if strings.Contains(authType, "Agent") {
			authIcon = "🎫"
		}

		info := []string{
			lipgloss.NewStyle().
				Foreground(accentCyan).
				Bold(true).
				Render("━━━ CONNECTION DETAILS ━━━"),
			"",
			labelStyle.Render("🏷️  Name:") + "    " + valueStyle.Render(conn.Name),
			labelStyle.Render("🌐 Host:") + "    " + valueStyle.Render(conn.Host),
			labelStyle.Render("👤 User:") + "    " + valueStyle.Render(conn.User),
			labelStyle.Render("🔌 Port:") + "    " + valueStyle.Render(fmt.Sprintf("%d", conn.Port)),
			labelStyle.Render(authIcon+" Auth:") + "    " + valueStyle.Render(authType),
			"",
			labelStyle.Render("🏷️  Tags:") + "    " + tagsStr.String(),
		}

		if conn.JumpHost != "" {
			info = append(info, "")
			info = append(info, labelStyle.Render("🔀 Jump:")+"    "+valueStyle.Render(conn.JumpHost))
		}

		info = append(info, "")
		info = append(info, labelStyle.Render("📊 Status:")+"   "+statusStyle.Render(" ✓ READY "))

		info = append(info, "")
		info = append(info, "")
		info = append(info, helpStyle.Render("Press Enter to connect • / to filter • q to quit"))

		dWidth := m.width - m.width/3 - 10

		if dWidth < 30 {
			dWidth = 30
		}
		dHeight := m.height - 8

		if dHeight < 10 {
			dHeight = 10
		}

		details = detailStyle.
			Width(dWidth).
			Height(dHeight).
			Render(strings.Join(info, "\n"))
	}

	listView := m.list.View()
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, details)

	// @todo: In fact brand name should come from the .env file
	header := headerStyle.Render("⚡ LEAP SSH MANAGER")
	subtitle := subtitleStyle.Render("Manage your SSH connections with ease")

	return appStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			subtitle,
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

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(primaryGreen).
		BorderForeground(primaryGreen).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(accentCyan).
		BorderForeground(primaryGreen)

	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = "📡 Connections"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(accentCyan).
		Bold(true).
		MarginLeft(1)
	l.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(primaryGreen).
		Bold(true)
	l.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(highlightColor)

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
