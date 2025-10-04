package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jatsandaruwan/logx/internal/config"
)

type UserManagementModel struct {
	cursor  int
	config  *config.Config
	users   []config.User
	mode    string // "menu", "list", "add", "delete"
	message string
}

func NewUserManagementMenu(cfg *config.Config) UserManagementModel {
	return UserManagementModel{
		config: cfg,
		users:  cfg.Users.Users,
		mode:   "menu",
	}
}

func (m UserManagementModel) Init() tea.Cmd {
	return nil
}

func (m UserManagementModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := 0
			switch m.mode {
			case "menu":
				maxCursor = 3 // 4 options
			case "list":
				maxCursor = len(m.users) - 1
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

func (m UserManagementModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.mode {
	case "menu":
		switch m.cursor {
		case 0: // Add User
			m.mode = "add"
			m.cursor = 0
		case 1: // List Users
			m.mode = "list"
			m.cursor = 0
		case 2: // Delete User
			m.mode = "delete"
			m.cursor = 0
		case 3: // Back
			mainMenu, _ := NewMainMenu()
			return mainMenu, nil
		}

	case "list":
		// Show user details (could expand this)
		if m.cursor < len(m.users) {
			user := m.users[m.cursor]
			m.message = fmt.Sprintf("User: %s | SSH Username: %s", user.Name, user.Username)
		}

	case "delete":
		if m.cursor < len(m.users) {
			// Delete user (simplified, should add confirmation)
			user := m.users[m.cursor]
			m.config.DeleteUser(user.Name)
			m.config.Save()
			m.message = fmt.Sprintf("✓ User '%s' deleted", user.Name)
			m.users = m.config.Users.Users
			m.mode = "menu"
			m.cursor = 0
		}
	}

	return m, nil
}

func (m UserManagementModel) View() string {
	var s strings.Builder

	// Title
	title := " 👤 User Management "
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

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
		s.WriteString(successStyle.Render(m.message))
	}

	// Help
	s.WriteString("\n\n")
	help := "↑/↓: Navigate • Enter: Select • Esc: Back"
	s.WriteString(helpStyle.Render(help))

	return s.String()
}

func (m UserManagementModel) renderMenu() string {
	options := []string{
		"➕ Add User",
		"📋 List Users",
		"🗑️  Delete User",
		"⬅️  Back to Main Menu",
	}

	var content string
	for i, option := range options {
		cursor := " "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶")
			option = focusedStyle.Render(option)
		} else {
			option = blurredStyle.Render(option)
		}
		content += fmt.Sprintf("%s %s\n", cursor, option)
	}

	return menuBoxStyle.Render(content)
}

func (m UserManagementModel) renderList() string {
	if len(m.users) == 0 {
		return errorStyle.Render("No users configured")
	}

	var content string
	for i, user := range m.users {
		cursor := " "
		line := fmt.Sprintf("%s (SSH: %s)", user.Name, user.Username)

		if i == m.cursor {
			cursor = cursorStyle.Render("▶")
			line = focusedStyle.Render(line)
		} else {
			line = blurredStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	return menuBoxStyle.Render(content)
}

func (m UserManagementModel) renderDelete() string {
	if len(m.users) == 0 {
		return successStyle.Render("No users to delete")
	}

	content := errorStyle.Render("⚠️  Select user to DELETE:\n\n")

	for i, user := range m.users {
		cursor := " "
		line := fmt.Sprintf("%s (SSH: %s)", user.Name, user.Username)

		if i == m.cursor {
			cursor = cursorStyle.Render("▶")
			line = focusedStyle.Render(line)
		} else {
			line = blurredStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	return menuBoxStyle.Render(content)
}
