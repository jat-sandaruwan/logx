package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jatsandaruwan/logx/internal/config"
)

// SettingsModel handles settings UI
type SettingsModel struct {
	cursor  int
	config  *config.Config
	mode    string // "menu", "editor"
	message string
	input   string
	width   int
	height  int
}

func NewSettingsMenu(cfg *config.Config) SettingsModel {
	return SettingsModel{
		config: cfg,
		mode:   "menu",
		width:  80,
		height: 24,
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.input = ""
			m.cursor = 0
			return m, nil

		case "up", "k":
			if m.mode == "menu" && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.mode == "menu" && m.cursor < 1 { // 2 options (0, 1)
				m.cursor++
			}

		case "enter":
			return m.handleSelection()

		case "backspace":
			if m.mode == "editor" && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			// Handle text input for editor mode
			if m.mode == "editor" && len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}

	return m, nil
}

func (m SettingsModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.mode {
	case "menu":
		switch m.cursor {
		case 0: // Set Editor
			m.mode = "editor"
			m.input = m.config.Editor
			if m.input == "" {
				m.input = ""
			}
			m.message = ""

		case 1: // Back
			mainMenu, _ := NewMainMenu()
			return mainMenu, nil
		}

	case "editor":
		// Save editor
		m.config.Editor = strings.TrimSpace(m.input)
		if err := m.config.Save(); err != nil {
			m.message = errorStyle.Render(fmt.Sprintf("❌ Error: %v", err))
		} else {
			editorMsg := "Platform Default"
			if m.config.Editor != "" {
				editorMsg = m.config.Editor
			}
			m.message = successStyle.Render(fmt.Sprintf("✓ Editor set to: %s", editorMsg))
		}
		m.mode = "menu"
		m.cursor = 0
		m.input = ""
	}

	return m, nil
}

func (m SettingsModel) View() string {
	var s strings.Builder

	// Title
	title := " ⚙️  Settings "
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	// Content based on mode
	switch m.mode {
	case "menu":
		s.WriteString(m.renderMenu())
	case "editor":
		s.WriteString(m.renderEditorInput())
	}

	// Message
	if m.message != "" {
		s.WriteString("\n\n")
		s.WriteString(m.message)
	}

	// Help
	s.WriteString("\n\n")
	if m.mode == "editor" {
		help := "Type editor command • Enter: Save • Esc: Cancel"
		s.WriteString(helpStyle.Render(help))
	} else {
		help := "↑/↓ or j/k: Navigate • Enter: Select • Esc: Back"
		s.WriteString(helpStyle.Render(help))
	}

	return s.String()
}

func (m SettingsModel) renderMenu() string {
	var content strings.Builder

	// Show current settings
	content.WriteString(labelStyle.Render("Current Settings"))
	content.WriteString("\n\n")

	currentEditor := m.config.Editor
	if currentEditor == "" {
		currentEditor = "Platform Default"
	}

	content.WriteString(blurredStyle.Render("✏️  Editor: "))
	content.WriteString(infoStyle.Render(currentEditor))
	content.WriteString("\n\n")

	// Menu options
	content.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	options := []string{
		"✏️  Change Editor",
		"⬅️  Back to Main Menu",
	}

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

func (m SettingsModel) renderEditorInput() string {
	var content strings.Builder

	content.WriteString(labelStyle.Render("Configure Editor"))
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(blurredStyle.Render("Enter the command to launch your editor:"))
	content.WriteString("\n\n")

	// Examples
	content.WriteString(infoStyle.Render("Examples:"))
	content.WriteString("\n")
	content.WriteString(blurredStyle.Render("  • code        (VS Code)"))
	content.WriteString("\n")
	content.WriteString(blurredStyle.Render("  • vim         (Vim)"))
	content.WriteString("\n")
	content.WriteString(blurredStyle.Render("  • nano        (Nano)"))
	content.WriteString("\n")
	content.WriteString(blurredStyle.Render("  • notepad++   (Notepad++)"))
	content.WriteString("\n")
	content.WriteString(blurredStyle.Render("  • subl        (Sublime Text)"))
	content.WriteString("\n\n")

	// Input field
	content.WriteString(labelStyle.Render("Editor Command:"))
	content.WriteString("\n")

	inputDisplay := m.input
	if inputDisplay == "" {
		inputDisplay = " "
	}
	inputDisplay += "█"

	content.WriteString(inputStyle.Render(inputDisplay))
	content.WriteString("\n\n")

	// Hint
	if m.input == "" {
		content.WriteString(helpStyle.Render("💡 Leave empty to use platform default"))
	} else {
		content.WriteString(blurredStyle.Render(fmt.Sprintf("Will use: %s", m.input)))
	}

	return menuBoxStyle.Render(content.String())
}
