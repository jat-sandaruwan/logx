package viewer

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1).
			MarginBottom(1)

	lineNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Width(6).
			Align(lipgloss.Right)

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	searchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFF00")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)
)

type LogViewerModel struct {
	content      []string
	serverName   string
	logFile      string
	width        int
	height       int
	offset       int
	cursor       int
	searchMode   bool
	searchQuery  string
	searchResult []int
	searchIndex  int
	message      string
}

func NewLogViewer(content []string, serverName, logFile string) LogViewerModel {
	return LogViewerModel{
		content:    content,
		serverName: serverName,
		logFile:    logFile,
		width:      80,
		height:     24,
		offset:     0,
		cursor:     0,
	}
}

func (m LogViewerModel) Init() tea.Cmd {
	return nil
}

func (m LogViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}

		case "down", "j":
			if m.cursor < len(m.content)-1 {
				m.cursor++
				visibleLines := m.height - 5
				if m.cursor >= m.offset+visibleLines {
					m.offset = m.cursor - visibleLines + 1
				}
			}

		case "pageup", "ctrl+u":
			m.cursor -= (m.height - 5) / 2
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.offset = m.cursor

		case "pagedown", "ctrl+d":
			m.cursor += (m.height - 5) / 2
			if m.cursor >= len(m.content) {
				m.cursor = len(m.content) - 1
			}
			visibleLines := m.height - 5
			if m.cursor >= m.offset+visibleLines {
				m.offset = m.cursor - visibleLines + 1
			}

		case "home", "g":
			m.cursor = 0
			m.offset = 0

		case "end", "G":
			m.cursor = len(m.content) - 1
			visibleLines := m.height - 5
			m.offset = m.cursor - visibleLines + 1
			if m.offset < 0 {
				m.offset = 0
			}

		case "/":
			m.searchMode = true
			m.searchQuery = ""
			m.searchResult = []int{}
			m.message = "Search: "

		case "n":
			if len(m.searchResult) > 0 {
				m.searchIndex = (m.searchIndex + 1) % len(m.searchResult)
				m.cursor = m.searchResult[m.searchIndex]
				m.ensureVisible()
			}

		case "N":
			if len(m.searchResult) > 0 {
				m.searchIndex--
				if m.searchIndex < 0 {
					m.searchIndex = len(m.searchResult) - 1
				}
				m.cursor = m.searchResult[m.searchIndex]
				m.ensureVisible()
			}

		case "s":
			return m, m.saveLog()
		}
	}

	return m, nil
}

func (m *LogViewerModel) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searchMode = false
		m.message = ""
		return m, nil

	case "enter":
		m.searchMode = false
		m.performSearch()
		if len(m.searchResult) > 0 {
			m.searchIndex = 0
			m.cursor = m.searchResult[0]
			m.ensureVisible()
			m.message = fmt.Sprintf("Found %d matches", len(m.searchResult))
		} else {
			m.message = "No matches found"
		}
		return m, nil

	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}

	m.message = "Search: " + m.searchQuery
	return m, nil
}

func (m *LogViewerModel) performSearch() {
	m.searchResult = []int{}
	query := strings.ToLower(m.searchQuery)

	for i, line := range m.content {
		if strings.Contains(strings.ToLower(line), query) {
			m.searchResult = append(m.searchResult, i)
		}
	}
}

func (m *LogViewerModel) ensureVisible() {
	visibleLines := m.height - 5
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visibleLines {
		m.offset = m.cursor - visibleLines + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m LogViewerModel) saveLog() tea.Cmd {
	return func() tea.Msg {
		filename := fmt.Sprintf("%s_%s.log", m.serverName, strings.ReplaceAll(m.logFile, "/", "_"))
		content := strings.Join(m.content, "\n")

		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error saving: %v", err))
		}

		return tea.Msg(fmt.Sprintf("Saved to: %s", filename))
	}
}

func (m LogViewerModel) View() string {
	var s strings.Builder

	// Title bar
	title := fmt.Sprintf(" 📋 %s - %s ", m.serverName, m.logFile)
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	// Content area
	visibleLines := m.height - 5
	start := m.offset
	end := m.offset + visibleLines
	if end > len(m.content) {
		end = len(m.content)
	}

	searchMap := make(map[int]bool)
	for _, idx := range m.searchResult {
		searchMap[idx] = true
	}

	for i := start; i < end; i++ {
		lineNum := lineNumberStyle.Render(fmt.Sprintf("%4d", i+1))
		line := m.content[i]

		// Highlight current line
		if i == m.cursor {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#2a2a2a")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Render(line)
		} else if searchMap[i] {
			// Highlight search results
			if m.searchQuery != "" {
				line = highlightSearch(line, m.searchQuery)
			}
		}

		s.WriteString(lineNum)
		s.WriteString(" ")
		s.WriteString(contentStyle.Render(line))
		s.WriteString("\n")
	}

	// Status bar
	s.WriteString("\n")
	status := fmt.Sprintf(" Line %d/%d ", m.cursor+1, len(m.content))
	if len(m.searchResult) > 0 {
		status += fmt.Sprintf("| Match %d/%d ", m.searchIndex+1, len(m.searchResult))
	}
	s.WriteString(statusStyle.Render(status))

	// Help bar or message
	s.WriteString("\n")
	if m.message != "" {
		s.WriteString(helpStyle.Render(m.message))
	} else {
		help := "↑↓: Navigate | /: Search | n/N: Next/Prev | s: Save | q: Quit"
		s.WriteString(helpStyle.Render(help))
	}

	return s.String()
}

func highlightSearch(line, query string) string {
	lower := strings.ToLower(line)
	queryLower := strings.ToLower(query)
	idx := strings.Index(lower, queryLower)

	if idx == -1 {
		return line
	}

	before := line[:idx]
	match := line[idx : idx+len(query)]
	after := line[idx+len(query):]

	return before + searchStyle.Render(match) + after
}

// OpenInternalViewer opens the log in the internal TUI viewer
func OpenInternalViewer(content []string, serverName, logFile string) error {
	m := NewLogViewer(content, serverName, logFile)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running viewer: %w", err)
	}

	return nil
}
