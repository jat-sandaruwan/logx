package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps SSH connection
type Client struct {
	conn *ssh.Client
}

// Connect establishes SSH connection
func Connect(host, username, password string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For production, use proper host key verification
		Timeout:         10 * time.Second,
	}

	// Add default port if not specified
	if !strings.Contains(host, ":") {
		host = host + ":22"
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	return &Client{conn: conn}, nil
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// FileExists checks if a file exists on the remote server
func (c *Client) FileExists(path string) (bool, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return false, err
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			_ = fmt.Errorf("an error occurred while closing the session %w", err)
		}
	}(session)

	cmd := fmt.Sprintf("test -f %s && echo exists", path)
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return false, nil
	}

	return strings.TrimSpace(string(output)) == "exists", nil
}

// DownloadFile downloads a file from remote server to local temp directory
func (c *Client) DownloadFile(remotePath string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			_ = fmt.Errorf("an error occurred while closing the session %w", err)
		}
	}(session)

	// Create temp file
	tmpFile, err := os.CreateTemp("", "logx-*.log")
	if err != nil {
		return "", err
	}
	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {
			_ = fmt.Errorf("an error occurred while closing the temp file %w", err)
		}
	}(tmpFile)

	// Use cat to read the file
	cmd := fmt.Sprintf("cat %s", remotePath)
	output, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := session.Start(cmd); err != nil {
		return "", err
	}

	// Copy to temp file
	if _, err := io.Copy(tmpFile, output); err != nil {
		return "", err
	}

	if err := session.Wait(); err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	return tmpFile.Name(), nil
}

// ListFiles lists files matching a pattern in a directory
func (c *Client) ListFiles(dir, pattern string) ([]string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			_ = fmt.Errorf("an error occurred while closing the session %w", err)
		}
	}(session)

	cmd := fmt.Sprintf("ls -1 %s 2>/dev/null | grep '%s' || true", dir, pattern)
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, f := range files {
		if f != "" {
			result = append(result, filepath.Join(dir, f))
		}
	}

	return result, nil
}
