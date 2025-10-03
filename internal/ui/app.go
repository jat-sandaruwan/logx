package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jatsandaruwan/logx/internal/config"
)

// AddAppInteractive shows interactive UI for adding an app
func AddAppInteractive() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Users.Users) == 0 {
		return fmt.Errorf("no users configured. Please add a user first using: logx user add")
	}

	var app config.App

	// App name
	fmt.Print("App Name: ")
	fmt.Scanln(&app.Name)

	// Select user
	fmt.Println("\nAvailable users:")
	for i, user := range cfg.Users.Users {
		fmt.Printf("  %d. %s (username: %s)\n", i+1, user.Name, user.Username)
	}
	fmt.Print("Select user (number): ")
	var userIdx int
	fmt.Scanln(&userIdx)
	if userIdx < 1 || userIdx > len(cfg.Users.Users) {
		return fmt.Errorf("invalid user selection")
	}
	app.UserRef = cfg.Users.Users[userIdx-1].ID

	// Log path
	fmt.Print("Log file path (e.g., /logs/testapp/testapp.log): ")
	fmt.Scanln(&app.LogPath)

	// Log pattern for rolling files
	fmt.Print("Log filename pattern with {date} placeholder (e.g., testapp-{date}.log or testapp.log-{date}): ")
	fmt.Scanln(&app.LogPattern)

	// Date format
	fmt.Println("\nDate format examples:")
	fmt.Println("  2006-01-02  -> 2025-09-10")
	fmt.Println("  20060102    -> 20250910")
	fmt.Println("  02-01-2006  -> 10-09-2025")
	fmt.Print("Date format (Go format): ")
	fmt.Scanln(&app.DateFormat)

	// Servers
	fmt.Println("\nEnter server IPs (one per line, empty line to finish):")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Server IP: ")
		server, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("error reading input: %w", err)
		}
		server = strings.TrimSpace(server)
		if server == "" {
			break
		}
		app.Servers = append(app.Servers, server)
	}

	if len(app.Servers) == 0 {
		return fmt.Errorf("at least one server is required")
	}

	// Save app
	if err := cfg.AddApp(app); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Printf("\n✓ App '%s' added successfully!\n", app.Name)
	return nil
}

// UpdateAppInteractive shows interactive UI for updating an app
func UpdateAppInteractive(appName string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	app, err := cfg.GetApp(appName)
	if err != nil {
		return err
	}

	fmt.Printf("Updating app: %s\n", app.Name)
	fmt.Println("Press Enter to keep current value, or type new value:")

	// Update fields
	fmt.Printf("User ref [%s]: ", app.UserRef)
	var input string
	fmt.Scanln(&input)
	if input != "" {
		app.UserRef = input
	}

	fmt.Printf("Log path [%s]: ", app.LogPath)
	input = ""
	fmt.Scanln(&input)
	if input != "" {
		app.LogPath = input
	}

	fmt.Printf("Log pattern [%s]: ", app.LogPattern)
	input = ""
	fmt.Scanln(&input)
	if input != "" {
		app.LogPattern = input
	}

	fmt.Printf("Date format [%s]: ", app.DateFormat)
	input = ""
	fmt.Scanln(&input)
	if input != "" {
		app.DateFormat = input
	}

	fmt.Println("\nCurrent servers:")
	for _, server := range app.Servers {
		fmt.Printf("  - %s\n", server)
	}
	fmt.Print("Update servers? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if strings.ToLower(input) == "y" {
		app.Servers = []string{}
		fmt.Println("Enter new servers (empty line to finish):")
		for {
			fmt.Print("Server IP: ")
			server, err := reader.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf("error reading input: %w", err)
			}
			server = strings.TrimSpace(server)
			if server == "" {
				break
			}
			app.Servers = append(app.Servers, server)
		}
	}

	if err := cfg.UpdateApp(*app); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Printf("\n✓ App '%s' updated successfully!\n", app.Name)
	return nil
}

// ListApps displays all configured apps
func ListApps() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Apps.Apps) == 0 {
		fmt.Println("No apps configured.")
		return nil
	}

	fmt.Println("Configured applications:")
	fmt.Println()
	for _, app := range cfg.Apps.Apps {
		fmt.Printf("Name: %s\n", app.Name)
		fmt.Printf("  User: %s\n", app.UserRef)
		fmt.Printf("  Path: %s\n", app.LogPath)
		fmt.Printf("  Pattern: %s\n", app.LogPattern)
		fmt.Printf("  Date Format: %s\n", app.DateFormat)
		fmt.Printf("  Servers: %s\n", strings.Join(app.Servers, ", "))
		fmt.Println()
	}

	return nil
}

// ListUsers displays all configured users
func ListUsers() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Users.Users) == 0 {
		fmt.Println("No users configured.")
		return nil
	}

	fmt.Println("Configured users:")
	fmt.Println()
	for _, user := range cfg.Users.Users {
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("  Username: %s\n", user.Username)
		fmt.Println()
	}

	return nil
}
