package ui

import (
	"fmt"
	"os"

	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/vault"
	"golang.org/x/term"
)

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
