package main

import (
	"fmt"
	"os"

	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/ui"
	"github.com/jatsandaruwan/logx/internal/vault"
)

const version = "1.0.0"

func main() {
	// If no arguments, show interactive TUI menu
	if len(os.Args) < 2 {
		if err := ui.RunMainMenu(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("logx version %s\n", version)

	case "help", "-h", "--help":
		printHelp()

	case "user":
		handleUserCommand()

	case "app":
		handleAppCommand()

	case "editor":
		handleEditorCommand()

	case "tui", "menu":
		// Explicit TUI mode
		if err := ui.RunMainMenu(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'logx' for interactive menu or 'logx help' for command list")
		os.Exit(1)
	}
}

func handleUserCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logx user <add|list|delete> [name]")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "add":
		if err := ui.AddUserInteractive(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := ui.ListUsers(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Usage: logx user delete <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		if err := deleteUser(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ User '%s' deleted successfully!\n", name)

	default:
		fmt.Printf("Unknown user subcommand: %s\n", subcommand)
		fmt.Println("Available: add, list, delete")
		os.Exit(1)
	}
}

func handleAppCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logx app <add|list|update|delete> [name]")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "add":
		if err := ui.AddAppInteractive(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := ui.ListApps(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: logx app update <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		if err := ui.UpdateAppInteractive(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Usage: logx app delete <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		if err := deleteApp(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ App '%s' deleted successfully!\n", name)

	default:
		fmt.Printf("Unknown app subcommand: %s\n", subcommand)
		fmt.Println("Available: add, list, update, delete")
		os.Exit(1)
	}
}

func handleEditorCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logx editor <set|show> [editor-command]")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "set":
		if len(os.Args) < 4 {
			fmt.Println("Usage: logx editor set <editor-command>")
			fmt.Println("Example: logx editor set \"code\"")
			os.Exit(1)
		}
		editor := os.Args[3]
		if err := setEditor(editor); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Editor set to: %s\n", editor)

	case "show":
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if cfg.Editor == "" {
			fmt.Println("No custom editor set. Using platform default.")
		} else {
			fmt.Printf("Current editor: %s\n", cfg.Editor)
		}

	default:
		fmt.Printf("Unknown editor subcommand: %s\n", subcommand)
		fmt.Println("Available: set, show")
		os.Exit(1)
	}
}

func deleteUser(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := cfg.DeleteUser(name); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	return vault.Delete(name)
}

func deleteApp(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := cfg.DeleteApp(name); err != nil {
		return err
	}

	return cfg.Save()
}

func setEditor(editor string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfg.Editor = editor
	return cfg.Save()
}

func printHelp() {
	fmt.Println("logx - Remote Log Viewer with Interactive TUI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  logx                           Launch interactive TUI menu")
	fmt.Println("  logx <command> [options]       Run command directly")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  tui, menu                      Launch interactive TUI menu")
	fmt.Println("  user <add|list|delete>         Manage users")
	fmt.Println("  app <add|list|update|delete>   Manage applications")
	fmt.Println("  editor <set|show>              Manage editor settings")
	fmt.Println("  version                        Show version")
	fmt.Println("  help                           Show this help")
	fmt.Println()
	fmt.Println("Interactive TUI Features:")
	fmt.Println("  • Colorful menu navigation with arrow keys")
	fmt.Println("  • Internal log viewer with search")
	fmt.Println("  • Save logs locally")
	fmt.Println("  • User and app management")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  logx                    # Launch interactive menu")
	fmt.Println("  logx user add           # Add user via CLI")
	fmt.Println("  logx app list           # List apps via CLI")
	fmt.Println()
	fmt.Println("Log Viewer Controls (in TUI):")
	fmt.Println("  ↑/↓ or j/k    Navigate lines")
	fmt.Println("  PgUp/PgDn     Page up/down")
	fmt.Println("  g/G           Go to top/bottom")
	fmt.Println("  /             Search")
	fmt.Println("  n/N           Next/previous match")
	fmt.Println("  s             Save log locally")
	fmt.Println("  q or Ctrl+C   Quit")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Config: ~/.config/logx/config.xml")
	fmt.Println("  Credentials: System keyring")
	fmt.Println()
}
