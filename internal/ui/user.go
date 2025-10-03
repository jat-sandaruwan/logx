package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/vault"
	"golang.org/x/term"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)

type userAddModel struct {
	step     int
	name     string
	username string
	password string
	err      error
	done     bool
}

func (m userAddModel) Init() tea.Cmd {
	return nil
}

func (m userAddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.step == 2 {
				// Save user
				if err := m.saveUser(); err != nil {
					m.err = err
					return m, nil
				}
				m.done = true
				return m, tea.Quit
			}
			m.step++
			return m, nil
		case "backspace":
			if m.step == 0 && len(m.name) > 0 {
				m.name = m.name[:len(m.name)-1]
			} else if m.step == 1 && len(m.username) > 0 {
				m.username = m.username[:len(m.username)-1]
			}
		default:
			if m.step == 0 && len(msg.String()) == 1 {
				m.name += msg.String()
			} else if m.step == 1 && len(msg.String()) == 1 {
				m.username += msg.String()
			}
		}
	}
	return m, nil
}

func (m userAddModel) View() string {
	if m.done {
		if m.err != nil {
			return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
		}
		return successStyle.Render(fmt.Sprintf("✓ User '%s' added successfully!\n", m.name))
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("Add New User"))
	s.WriteString("\n\n")

	if m.step == 0 {
		s.WriteString(labelStyle.Render("User Name (identifier): "))
		s.WriteString(inputStyle.Render(m.name))
		s.WriteString(inputStyle.Render("█"))
	} else {
		s.WriteString(labelStyle.Render("User Name: "))
		s.WriteString(inputStyle.Render(m.name))
		s.WriteString("\n\n")
	}

	if m.step >= 1 {
		if m.step == 1 {
			s.WriteString(labelStyle.Render("SSH Username: "))
			s.WriteString(inputStyle.Render(m.username))
			s.WriteString(inputStyle.Render("█"))
		} else {
			s.WriteString(labelStyle.Render("SSH Username: "))
			s.WriteString(inputStyle.Render(m.username))
			s.WriteString("\n\n")
		}
	}

	if m.step >= 2 {
		s.WriteString(labelStyle.Render("SSH Password: "))
		s.WriteString(inputStyle.Render(strings.Repeat("*", len(m.password))))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("Press Enter to save"))
	}

	if m.err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return s.String()
}

func (m *userAddModel) saveUser() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	user := config.User{
		ID:       m.name,
		Name:     m.name,
		Username: m.username,
	}

	if err := cfg.AddUser(user); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	// Store credentials in vault
	creds := vault.Credentials{
		Username: m.username,
		Password: m.password,
	}
	return vault.Store(m.name, creds)
}

// AddUserInteractive shows interactive UI for adding a user
func AddUserInteractive() error {
	// Read password from terminal
	fmt.Print("User Name (identifier): ")
	var name string
	fmt.Scanln(&name)

	fmt.Print("SSH Username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("SSH Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	password := string(passwordBytes)
	fmt.Println()

	// Save user
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	user := config.User{
		ID:       name,
		Name:     name,
		Username: username,
	}

	if err := cfg.AddUser(user); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	// Store credentials
	creds := vault.Credentials{
		Username: username,
		Password: password,
	}
	if err := vault.Store(name, creds); err != nil {
		return err
	}

	fmt.Printf("✓ User '%s' added successfully!\n", name)
	return nil
}