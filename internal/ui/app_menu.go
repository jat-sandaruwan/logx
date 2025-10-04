package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jatsandaruwan/logx/internal/config"
)

// AppManagementModel handles app management UI
type AppManagementModel struct {
	cursor  int
	config  *config.Config
	apps    []config.App
	mode    string // "menu", "list", "delete"
	message string
	width   int
	height  int
}

func NewAppManagementMenu(cfg *config.Config) AppManagementModel {
	return AppManagementModel{
		config: cfg,
		apps:   cfg.Apps.Apps,
		mode:   "menu",
		width:  80,
		height: 24,
	}
}

func (m AppManagementModel) Init() tea.Cmd {
	return nil
}

func (m AppManagementModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.mode == "menu" {
				// Go back to main menu
				mainMenu, _ := NewMainMenu()
				return mainMenu, nil
			}
			m.mode = "menu"
			m.message = ""
			m.cursor = 0
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := 0
			switch m.mode {
			case "menu":
				maxCursor = 2 // 3 options (0, 1, 2)
			case "list", "delete":
				if len(m.apps) > 0 {
					maxCursor = len(m.apps) - 1
				}
			}
			if m.cursor < maxCursor {
				m.cursor++
			}

		case "enter":
			return m.handleSelection()
		}
	}

	return m, nil
}

func (m AppManagementModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.mode {
	case "menu":
		switch m.cursor {
		case 0: // List Apps
			m.mode = "list"
			m.cursor = 0
		case 1: // Delete App
			if len(m.apps) == 0 {
				m.message = errorStyle.Render("No apps to delete")
			} else {
				m.mode = "delete"
				m.cursor = 0
			}
		case 2: // Back
			mainMenu, _ := NewMainMenu()
			return mainMenu, nil
		}

	case "list":
		if m.cursor < len(m.apps) {
			app := m.apps[m.cursor]
			m.message = infoStyle.Render(
				fmt.Sprintf("📱 %s | 🖥️  %d servers | 📝 Pattern: %s | 📅 Format: %s",
					app.Name, len(app.Servers), app.LogPattern, app.DateFormat))
		}

	case "delete":
		if m.cursor < len(m.apps) {
			app := m.apps[m.cursor]
			if err := m.config.DeleteApp(app.Name); err != nil {
				m.message = errorStyle.Render(fmt.Sprintf("❌ Error: %v", err))
			} else {
				if err := m.config.Save(); err != nil {
					m.message = errorStyle.Render(fmt.Sprintf("❌ Error saving: %v", err))
				} else {
					m.message = successStyle.Render(fmt.Sprintf("✓ App '%s' deleted successfully!", app.Name))
					m.apps = m.config.Apps.Apps
					m.mode = "menu"
					m.cursor = 0
				}
			}
		}
	}

	return m, nil
}

func (m AppManagementModel) View() string {
	var s strings.Builder

	// Title
	title := " 📱 App Management "
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	// Content based on mode
	switch m.mode {
	case "menu":
		s.WriteString(m.renderMenu())
	case "list":
		s.WriteString(m.renderList())
	case "delete":
		s.WriteString(m.renderDelete())
	}

	// Message
	if m.message != "" {
		s.WriteString("\n\n")
		s.WriteString(m.message)
	}

	// Help
	s.WriteString("\n\n")
	help := "↑/↓ or j/k: Navigate • Enter: Select • Esc: Back"
	s.WriteString(helpStyle.Render(help))

	return s.String()
}

func (m AppManagementModel) renderMenu() string {
	options := []string{
		"📋 List Apps",
		"🗑️  Delete App",
		"⬅️  Back to Main Menu",
	}

	var content strings.Builder
	content.WriteString("Select an option:\n\n")

	for i, option := range options {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
			option = focusedStyle.Render(option)
		} else {
			option = blurredStyle.Render(option)
		}
		content.WriteString(fmt.Sprintf("%s%s\n", cursor, option))
	}

	return menuBoxStyle.Render(content.String())
}

func (m AppManagementModel) renderList() string {
	var content strings.Builder

	if len(m.apps) == 0 {
		content.WriteString(errorStyle.Render("❌ No apps configured"))
		content.WriteString("\n\n")
		content.WriteString(blurredStyle.Render("Add apps using: logx app add"))
		return menuBoxStyle.Render(content.String())
	}

	content.WriteString(fmt.Sprintf("Configured Applications (%d):\n\n", len(m.apps)))

	for i, app := range m.apps {
		cursor := "  "
		line := fmt.Sprintf("📱 %s (%d servers)", app.Name, len(app.Servers))

		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
			line = focusedStyle.Render(line)
		} else {
			line = blurredStyle.Render(line)
		}
		content.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("💡 Press Enter to view app details"))

	return menuBoxStyle.Render(content.String())
}

func (m AppManagementModel) renderDelete() string {
	var content strings.Builder

	if len(m.apps) == 0 {
		content.WriteString(errorStyle.Render("❌ No apps to delete"))
		return menuBoxStyle.Render(content.String())
	}

	content.WriteString(warningStyle.Render("⚠️  WARNING: Select app to DELETE"))
	content.WriteString("\n\n")
	content.WriteString(blurredStyle.Render("This action cannot be undone!"))
	content.WriteString("\n\n")

	for i, app := range m.apps {
		cursor := "  "
		line := fmt.Sprintf("📱 %s", app.Name)

		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
			line = focusedStyle.Render(line)
		} else {
			line = blurredStyle.Render(line)
		}
		content.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	return menuBoxStyle.Render(content.String())
}
