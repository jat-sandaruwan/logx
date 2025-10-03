package vault

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "logx"

// Credentials holds SSH credentials
type Credentials struct {
	Username string
	Password string
}

// Store saves credentials to system keyring
func Store(userID string, creds Credentials) error {
	// Store in format: username:password
	secret := fmt.Sprintf("%s:%s", creds.Username, creds.Password)
	return keyring.Set(serviceName, userID, secret)
}

// Get retrieves credentials from system keyring
func Get(userID string) (*Credentials, error) {
	secret, err := keyring.Get(serviceName, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	// Parse username:password format
	var username, password string
	_, err = fmt.Sscanf(secret, "%s:%s", &username, &password)
	if err != nil {
		// Try to parse with spaces or special chars
		for i := 0; i < len(secret); i++ {
			if secret[i] == ':' {
				username = secret[:i]
				if i+1 < len(secret) {
					password = secret[i+1:]
				}
				break
			}
		}
		if username == "" {
			return nil, fmt.Errorf("invalid credential format")
		}
	}

	return &Credentials{
		Username: username,
		Password: password,
	}, nil
}

// Delete removes credentials from system keyring
func Delete(userID string) error {
	return keyring.Delete(serviceName, userID)
}

// Exists checks if credentials exist for a user
func Exists(userID string) bool {
	_, err := keyring.Get(serviceName, userID)
	return err == nil
}
