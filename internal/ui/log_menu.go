package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/ssh"
	"github.com/jatsandaruwan/logx/internal/vault"
	"github.com/jatsandaruwan/logx/internal/viewer"
)

type LogSelectionModel struct {
	cursor      int
	config      *config.Config
	apps        []config.App
	mode        string // "select", "date", "server", "loading", "view"
	selectedApp *config.App
	dateInput   string
	servers     []string
	serverIdx   int
	loading     bool
	message     string
	logContent  []string
}

func NewLogSelectionMenu(cfg *config.Config) LogSelectionModel {
	return LogSelectionModel{
		config: cfg,
		apps:   cfg.Apps.Apps,
		mode:   "select",
	}
}

func (m LogSelectionModel) Init() tea.Cmd {
	return nil
}

func (m LogSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.mode == "select" {
				mainMenu, _ := NewMainMenu()
				return mainMenu, nil
			}
			m.mode = "select"
			m.cursor = 0
			m.message = ""
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := 0
			switch m.mode {
			case "select":
				maxCursor = len(m.apps)
			case "server":
				maxCursor = len(m.servers)
			}
			if m.cursor < maxCursor {
				m.cursor++
			}

		case "enter":
			return m.handleSelection()

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "-":
			if m.mode == "date" {
				m.dateInput += msg.String()
			}

		case "backspace":
			if m.mode == "date" && len(m.dateInput) > 0 {
				m.dateInput = m.dateInput[:len(m.dateInput)-1]
			}
		}

	case loadingMsg:
		m.loading = false
		if msg.err != nil {
			m.message = errorStyle.Render(fmt.Sprintf("Error: %v", msg.err))
			m.mode = "select"
		} else {
			m.logContent = msg.content
			m.mode = "view"
			// Launch internal viewer
			return m, func() tea.Msg {
				viewer.OpenInternalViewer(msg.content, msg.server, msg.logFile)
				return backToMenuMsg{}
			}
		}

	case backToMenuMsg:
		m.mode = "select"
		m.cursor = 0
	}

	return m, nil
}

type loadingMsg struct {
	content []string
	server  string
	logFile string
	err     error
}

type backToMenuMsg struct{}

func (m LogSelectionModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.mode {
	case "select":
		if m.cursor == len(m.apps) {
			// Back option
			mainMenu, _ := NewMainMenu()
			return mainMenu, nil
		}
		m.selectedApp = &m.apps[m.cursor]
		m.servers = append([]string{"All Servers"}, m.selectedApp.Servers...)
		m.mode = "server"
		m.cursor = 0

	case "server":
		if m.cursor == 0 {
			// All servers - just pick first for now
			m.serverIdx = 0
		} else {
			m.serverIdx = m.cursor - 1
		}
		m.mode = "date"
		m.dateInput = time.Now().Format("2006-01-02")

	case "date":
		// Load logs
		m.loading = true
		return m, m.loadLogs()
	}

	return m, nil
}

func (m LogSelectionModel) loadLogs() tea.Cmd {
	return func() tea.Msg {
		// Get user credentials
		user, err := m.config.GetUser(m.selectedApp.UserRef)
		if err != nil {
			return loadingMsg{err: err}
		}

		creds, err := vault.Get(user.ID)
		if err != nil {
			return loadingMsg{err: err}
		}

		// Parse date
		logDate, err := time.Parse("2006-01-02", m.dateInput)
		if err != nil {
			return loadingMsg{err: fmt.Errorf("invalid date format")}
		}

		// Format the log filename
		formattedDate := logDate.Format(m.selectedApp.DateFormat)
		logFileName := strings.ReplaceAll(m.selectedApp.LogPattern, "{date}", formattedDate)
		logFilePath := m.selectedApp.LogPath
		if strings.Contains(logFilePath, "/") {
			logFilePath = logFilePath[:strings.LastIndex(logFilePath, "/")+1] + logFileName
		}

		// Connect to server
		server := m.selectedApp.Servers[m.serverIdx]
		client, err := ssh.Connect(server, creds.Username, creds.Password)
		if err != nil {
			return loadingMsg{err: fmt.Errorf("failed to connect to %s: %w", server, err)}
		}
		defer client.Close()

		// Check if file exists
		exists, err := client.FileExists(logFilePath)
		if err != nil || !exists {
			return loadingMsg{err: fmt.Errorf("log file not found: %s", logFilePath)}
		}

		// Download and read file
		localPath, err := client.DownloadFile(logFilePath)
		if err != nil {
			return loadingMsg{err: fmt.Errorf("failed to download: %w", err)}
		}

		// Read content
		contentBytes, err := os.ReadFile(localPath)
		if err != nil {
			return loadingMsg{err: fmt.Errorf("failed to read file: %w", err)}
		}

		lines := strings.Split(string(contentBytes), "\n")

		return loadingMsg{
			content: lines,
			server:  server,
			logFile: logFileName,
		}
	}
}

var (
	// Main Menu Styles
	logFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 2)

	logBlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 2)

	logTitleBannerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Bold(true).
				Padding(1, 4).
				MarginBottom(2)

	logMenuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(50)

	logStatsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	logCursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	logHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)

	logServerTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFA500")).
				Bold(true)
)

func (m LogSelectionModel) View() string {
	var s strings.Builder

	// Title
	title := " 📋 Log Viewer "
	s.WriteString(logTitleBannerStyle.Render(title))
	s.WriteString("\n\n")

	if m.loading {
		s.WriteString(logStatsStyle.Render("⏳ Loading logs..."))
		return s.String()
	}

	switch m.mode {
	case "select":
		s.WriteString(m.renderAppSelect())
	case "server":
		s.WriteString(m.renderServerSelect())
	case "date":
		s.WriteString(m.renderDateInput())
	}

	// Message
	if m.message != "" {
		s.WriteString("\n\n")
		s.WriteString(m.message)
	}

	// Help
	s.WriteString("\n\n")
	if m.mode == "date" {
		help := "Type date (YYYY-MM-DD) • Enter: View • Esc: Back"
		s.WriteString(logHelpStyle.Render(help))
	} else {
		help := "↑/↓: Navigate • Enter: Select • Esc: Back"
		s.WriteString(logHelpStyle.Render(help))
	}

	return s.String()
}

func (m LogSelectionModel) renderAppSelect() string {
	if len(m.apps) == 0 {
		return errorStyle.Render("No apps configured")
	}

	content := "Select an application:\n\n"

	for i, app := range m.apps {
		cursor := " "
		line := fmt.Sprintf("📱 %s (%d servers)", app.Name, len(app.Servers))

		if i == m.cursor {
			cursor = logCursorStyle.Render("▶")
			line = logFocusedStyle.Render(line)
		} else {
			line = logBlurredStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	// Back option
	cursor := " "
	line := "⬅️  Back to Main Menu"
	if m.cursor == len(m.apps) {
		cursor = logCursorStyle.Render("▶")
		line = logFocusedStyle.Render(line)
	} else {
		line = logBlurredStyle.Render(line)
	}
	content += fmt.Sprintf("\n%s %s\n", cursor, line)

	return logMenuBoxStyle.Render(content)
}

func (m LogSelectionModel) renderServerSelect() string {
	content := fmt.Sprintf("Select server for %s:\n\n",
		logFocusedStyle.Render(m.selectedApp.Name))

	for i, server := range m.servers {
		cursor := " "
		line := server
		if i == 0 {
			line = "🌐 " + line
		} else {
			line = "🖥️  " + line
		}

		if i == m.cursor {
			cursor = logCursorStyle.Render("▶")
			line = logFocusedStyle.Render(line)
		} else {
			line = logBlurredStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	return logMenuBoxStyle.Render(content)
}

func (m LogSelectionModel) renderDateInput() string {
	content := fmt.Sprintf("Viewing logs for: %s\n",
		logFocusedStyle.Render(m.selectedApp.Name))

	serverName := "All Servers"
	if m.serverIdx >= 0 && m.serverIdx < len(m.selectedApp.Servers) {
		serverName = m.selectedApp.Servers[m.serverIdx]
	}
	content += fmt.Sprintf("Server: %s\n\n",
		logServerTagStyle.Render(serverName))

	content += "Enter date (YYYY-MM-DD):\n"
	content += logFocusedStyle.Render(m.dateInput + "█")
	content += "\n\n"
	content += logBlurredStyle.Render("Press Enter to view logs")

	return logBlurredStyle.Render(content)
}
