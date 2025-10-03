package viewer

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jatsandaruwan/logx/internal/config"
	"github.com/jatsandaruwan/logx/internal/editor"
	"github.com/jatsandaruwan/logx/internal/ssh"
	"github.com/jatsandaruwan/logx/internal/vault"
)

// ViewLogs opens log files for the specified app and date
func ViewLogs(appName, dateStr, serverFilter string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	app, err := cfg.GetApp(appName)
	if err != nil {
		return err
	}

	// Get user credentials
	user, err := cfg.GetUser(app.UserRef)
	if err != nil {
		return err
	}

	creds, err := vault.Get(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get credentials for user %s: %w", user.Name, err)
	}

	// Parse date
	var logDate time.Time
	if dateStr == "" {
		logDate = time.Now()
	} else {
		logDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format. Use YYYY-MM-DD: %w", err)
		}
	}

	// Format the log filename
	formattedDate := logDate.Format(app.DateFormat)
	logFileName := strings.ReplaceAll(app.LogPattern, "{date}", formattedDate)
	logDir := filepath.Dir(app.LogPath)
	logFilePath := filepath.Join(logDir, logFileName)

	fmt.Printf("Looking for logs: %s\n", logFileName)
	fmt.Printf("Date: %s\n\n", logDate.Format("2006-01-02"))

	// Filter servers if specified
	servers := app.Servers
	if serverFilter != "" {
		servers = []string{serverFilter}
	}

	var downloadedFiles []string

	// Connect to each server and download logs
	for _, server := range servers {
		fmt.Printf("Connecting to %s...\n", server)

		client, err := ssh.Connect(server, creds.Username, creds.Password)
		if err != nil {
			fmt.Printf("  ✗ Failed to connect: %v\n", err)
			continue
		}

		// Check if file exists
		exists, err := client.FileExists(logFilePath)
		if err != nil {
			fmt.Printf("  ✗ Error checking file: %v\n", err)
			client.Close()
			continue
		}

		if !exists {
			fmt.Printf("  ✗ Log file not found: %s\n", logFilePath)
			client.Close()
			continue
		}

		// Download file
		fmt.Printf("  ↓ Downloading log file...\n")
		localPath, err := client.DownloadFile(logFilePath)
		if err != nil {
			fmt.Printf("  ✗ Failed to download: %v\n", err)
			client.Close()
			continue
		}

		fmt.Printf("  ✓ Downloaded to: %s\n", localPath)
		downloadedFiles = append(downloadedFiles, localPath)

		client.Close()
	}

	if len(downloadedFiles) == 0 {
		return fmt.Errorf("no log files found for the specified date")
	}

	// Open each file in editor
	fmt.Println("\nOpening log files...")
	for _, file := range downloadedFiles {
		if cfg.Editor != "" {
			if err := editor.OpenWithCustom(file, cfg.Editor); err != nil {
				fmt.Printf("Failed to open %s: %v\n", file, err)
			}
		} else {
			if err := editor.Open(file); err != nil {
				fmt.Printf("Failed to open %s: %v\n", file, err)
			}
		}
	}

	return nil
}

// ViewCurrentLogs opens the current (non-dated) log file
func ViewCurrentLogs(appName, serverFilter string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	app, err := cfg.GetApp(appName)
	if err != nil {
		return err
	}

	// Get user credentials
	user, err := cfg.GetUser(app.UserRef)
	if err != nil {
		return err
	}

	creds, err := vault.Get(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get credentials for user %s: %w", user.Name, err)
	}

	fmt.Printf("Looking for current logs: %s\n\n", app.LogPath)

	// Filter servers if specified
	servers := app.Servers
	if serverFilter != "" {
		servers = []string{serverFilter}
	}

	var downloadedFiles []string

	// Connect to each server and download logs
	for _, server := range servers {
		fmt.Printf("Connecting to %s...\n", server)

		client, err := ssh.Connect(server, creds.Username, creds.Password)
		if err != nil {
			fmt.Printf("  ✗ Failed to connect: %v\n", err)
			continue
		}

		// Check if file exists
		exists, err := client.FileExists(app.LogPath)
		if err != nil {
			fmt.Printf("  ✗ Error checking file: %v\n", err)
			client.Close()
			continue
		}

		if !exists {
			fmt.Printf("  ✗ Log file not found: %s\n", app.LogPath)
			client.Close()
			continue
		}

		// Download file
		fmt.Printf("  ↓ Downloading log file...\n")
		localPath, err := client.DownloadFile(app.LogPath)
		if err != nil {
			fmt.Printf("  ✗ Failed to download: %v\n", err)
			client.Close()
			continue
		}

		fmt.Printf("  ✓ Downloaded to: %s\n", localPath)
		downloadedFiles = append(downloadedFiles, localPath)

		client.Close()
	}

	if len(downloadedFiles) == 0 {
		return fmt.Errorf("no log files found")
	}

	// Open each file in editor
	fmt.Println("\nOpening log files...")
	for _, file := range downloadedFiles {
		if cfg.Editor != "" {
			if err := editor.OpenWithCustom(file, cfg.Editor); err != nil {
				fmt.Printf("Failed to open %s: %v\n", file, err)
			}
		} else {
			if err := editor.Open(file); err != nil {
				fmt.Printf("Failed to open %s: %v\n", file, err)
			}
		}
	}

	return nil
}