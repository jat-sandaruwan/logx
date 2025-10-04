package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jatsandaruwan/logx/internal/config"
)

type MainMenuModel struct {
	cursor    int
	options   []string
	config    *config.Config
	userCount int
	appCount  int
	quitting  bool
	width     int
	height    int
}

func NewMainMenu() (*MainMenuModel, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &MainMenuModel{
		options: []string{
			"👤 User Management",
			"📱 App Management",
			"📋 View Logs",
			"⚙️ Settings",
			"❌ Exit",
		},
		config:    cfg,
		userCount: len(cfg.Users.Users),
		appCount:  len(cfg.Apps.Apps),
		width:     80,
		height:    24,
	}, nil
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter":
			return m.handleSelection()
		}
	}

	return m, nil
}

func (m MainMenuModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0: // User Management
		return NewUserManagementMenu(m.config), nil
	case 1: // App Management
		return NewAppManagementMenu(m.config), nil
	case 2: // View Logs
		return NewLogSelectionMenu(m.config), nil
	case 3: // Settings
		return NewSettingsMenu(m.config), nil
	case 4: // Exit
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	if m.quitting {
		return statusStyle.Render("\n✓ Goodbye!\n\n")
	}

	var s strings.Builder

	// Banner
	banner := ` 
 ██╗      ██████╗  ██████╗ ██╗  ██╗
 ██║     ██╔═══██╗██╔════╝ ╚██╗██╔╝
 ██║     ██║   ██║██║  ███╗ ╚███╔╝ 
 ██║     ██║   ██║██║   ██║ ██╔██╗ 
 ███████╗╚██████╔╝╚██████╔╝██╔╝ ██╗
 ╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝
    Remote Log Viewer v1.0.0
`
	s.WriteString(titleBannerStyle.Render(banner))
	s.WriteString("\n\n")

	// Stats
	stats := fmt.Sprintf("👥 %d Users  •  📱 %d Apps", m.userCount, m.appCount)
	s.WriteString(statusStyle.Render(stats))
	s.WriteString("\n\n")

	// Menu
	var menuContent strings.Builder
	for i, option := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
			option = focusedStyle.Render(option)
		} else {
			option = blurredStyle.Render(option)
		}
		menuContent.WriteString(fmt.Sprintf("%s%s\n", cursor, option))
	}

	s.WriteString(menuBoxStyle.Render(menuContent.String()))
	s.WriteString("\n\n")

	// Help
	help := "↑/↓ or j/k: Navigate • Enter: Select • q: Quit"
	s.WriteString(helpStyle.Render(help))

	return s.String()
}

// RunMainMenu starts the main menu TUI
func RunMainMenu() error {
	m, err := NewMainMenu()
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
