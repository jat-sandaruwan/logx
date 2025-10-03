package config

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the root configuration
type Config struct {
	XMLName xml.Name `xml:"config"`
	Users   Users    `xml:"users"`
	Apps    Apps     `xml:"apps"`
	Editor  string   `xml:"editor,omitempty"`
}

// Users contains all user configurations
type Users struct {
	Users []User `xml:"user"`
}

// User represents SSH user credentials reference
type User struct {
	ID       string `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Username string `xml:"username,attr"`
}

// Apps contains all application configurations
type Apps struct {
	Apps []App `xml:"app"`
}

// App represents an application configuration
type App struct {
	Name       string   `xml:"name,attr"`
	UserRef    string   `xml:"user-ref"`
	LogPath    string   `xml:"log-path"`
	LogPattern string   `xml:"log-pattern"`
	DateFormat string   `xml:"date-format"`
	Servers    []string `xml:"servers>server"`
}

// GetConfigPath returns the platform-specific config file path
func GetConfigPath() (string, error) {
	var configDir string

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		configDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}

	logxDir := filepath.Join(configDir, "logx")
	if err := os.MkdirAll(logxDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(logxDir, "config.xml"), nil
}

// Load reads the configuration from XML file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Users: Users{Users: []User{}},
			Apps:  Apps{Apps: []App{}},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save writes the configuration to XML file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Add XML header
	xmlData := []byte(xml.Header + string(data))

	return os.WriteFile(configPath, xmlData, 0600)
}

// GetUser finds a user by ID
func (c *Config) GetUser(id string) (*User, error) {
	for _, user := range c.Users.Users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found: %s", id)
}

// GetUserByName finds a user by name
func (c *Config) GetUserByName(name string) (*User, error) {
	for _, user := range c.Users.Users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found: %s", name)
}

// AddUser adds a new user to configuration
func (c *Config) AddUser(user User) error {
	// Check if user already exists
	for _, u := range c.Users.Users {
		if u.ID == user.ID || u.Name == user.Name {
			return fmt.Errorf("user already exists: %s", user.Name)
		}
	}
	c.Users.Users = append(c.Users.Users, user)
	return nil
}

// DeleteUser removes a user from configuration
func (c *Config) DeleteUser(name string) error {
	for i, user := range c.Users.Users {
		if user.Name == name {
			c.Users.Users = append(c.Users.Users[:i], c.Users.Users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user not found: %s", name)
}

// GetApp finds an app by name
func (c *Config) GetApp(name string) (*App, error) {
	for _, app := range c.Apps.Apps {
		if app.Name == name {
			return &app, nil
		}
	}
	return nil, fmt.Errorf("app not found: %s", name)
}

// AddApp adds a new app to configuration
func (c *Config) AddApp(app App) error {
	// Check if app already exists
	for _, a := range c.Apps.Apps {
		if a.Name == app.Name {
			return fmt.Errorf("app already exists: %s", app.Name)
		}
	}
	c.Apps.Apps = append(c.Apps.Apps, app)
	return nil
}

// UpdateApp updates an existing app
func (c *Config) UpdateApp(app App) error {
	for i, a := range c.Apps.Apps {
		if a.Name == app.Name {
			c.Apps.Apps[i] = app
			return nil
		}
	}
	return fmt.Errorf("app not found: %s", app.Name)
}

// DeleteApp removes an app from configuration
func (c *Config) DeleteApp(name string) error {
	for i, app := range c.Apps.Apps {
		if app.Name == name {
			c.Apps.Apps = append(c.Apps.Apps[:i], c.Apps.Apps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("app not found: %s", name)
}
