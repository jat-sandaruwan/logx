package main

import (
	"fmt"
	"os"

	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/ui"
	"github.com/jatsandaruwan/logx/internal/vault"
	"github.com/jatsandaruwan/logx/internal/viewer"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
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

	default:
		// Assume it's an app name for viewing logs
		handleViewLogs()
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

func handleViewLogs() {
	appName := os.Args[1]
	var dateStr string
	var serverFilter string

	// Parse flags
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--server" || arg == "-s" {
			if i+1 < len(os.Args) {
				serverFilter = os.Args[i+1]
				i++
			}
		} else if dateStr == "" {
			dateStr = arg
		}
	}

	var err error
	if dateStr == "" {
		// View current log
		err = viewer.ViewCurrentLogs(appName, serverFilter)
	} else {
		// View dated log
		err = viewer.ViewLogs(appName, dateStr, serverFilter)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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

	// Delete from vault
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

func printUsage() {
	fmt.Println("Usage: logx <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  user <add|list|delete>           Manage users")
	fmt.Println("  app <add|list|update|delete>     Manage applications")
	fmt.Println("  editor <set|show>                Manage editor settings")
	fmt.Println("  <appname> [date] [--server IP]   View logs")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  logx user add")
	fmt.Println("  logx app add")
	fmt.Println("  logx app list")
	fmt.Println("  logx testapp                     View current logs")
	fmt.Println("  logx testapp 2025-09-10          View logs for specific date")
	fmt.Println("  logx testapp --server 192.168.0.1  View logs from specific server")
	fmt.Println()
	fmt.Println("Run 'logx help' for more information.")
}

func printHelp() {
	fmt.Println("logx - Remote Log Viewer")
	fmt.Println()
	printUsage()
	fmt.Println()
	fmt.Println("User Management:")
	fmt.Println("  logx user add              Add a new SSH user with credentials")
	fmt.Println("  logx user list             List all configured users")
	fmt.Println("  logx user delete <name>    Delete a user and credentials")
	fmt.Println()
	fmt.Println("Application Management:")
	fmt.Println("  logx app add               Add a new application configuration")
	fmt.Println("  logx app list              List all configured applications")
	fmt.Println("  logx app update <name>     Update an application configuration")
	fmt.Println("  logx app delete <name>     Delete an application configuration")
	fmt.Println()
	fmt.Println("Editor Management:")
	fmt.Println("  logx editor set <cmd>      Set custom editor command")
	fmt.Println("  logx editor show           Show current editor setting")
	fmt.Println()
	fmt.Println("View Logs:")
	fmt.Println("  logx <appname>                    View current log files")
	fmt.Println("  logx <appname> <YYYY-MM-DD>       View logs for specific date")
	fmt.Println("  logx <appname> --server <IP>      View logs from specific server only")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Config file: ~/.config/logx/config.xml")
	fmt.Println("  Credentials: Stored in system keyring")
	fmt.Println()
}