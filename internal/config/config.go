package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/paramientos/leap/internal/utils"
	"gopkg.in/yaml.v3"
)

type Connection struct {
	Name         string    `yaml:"name"`
	Host         string    `yaml:"host"`
	User         string    `yaml:"user"`
	Port         int       `yaml:"port"`
	Password     string    `yaml:"password,omitempty"`
	IdentityFile string    `yaml:"identity_file,omitempty"`
	Tags         []string  `yaml:"tags,omitempty"`
	JumpHost     string    `yaml:"jump_host,omitempty"`
	Tunnels      []Tunnel  `yaml:"tunnels,omitempty"`
	LastUsed     time.Time `yaml:"last_used,omitempty"`
	Favorite     bool      `yaml:"favorite,omitempty"`
	Notes        string    `yaml:"notes,omitempty"`
	UsageCount   int       `yaml:"usage_count,omitempty"`
	CreatedAt    time.Time `yaml:"created_at,omitempty"`
}

type Tunnel struct {
	Local  int `yaml:"local"`
	Remote int `yaml:"remote"`
}

type Config struct {
	Connections map[string]Connection `yaml:"connections"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".leap", "connections.yaml")
}

func LoadConfig(passphrase string) (*Config, error) {
	path := GetConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Connections: make(map[string]Connection)}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Check if data is encrypted with age
	isEncrypted := bytes.HasPrefix(data, []byte("age-encryption.org"))

	if isEncrypted {
		if passphrase == "" {
			return nil, fmt.Errorf("config is encrypted but LEAP_MASTER_PASSWORD is not set")
		}
		decrypted, err := utils.Decrypt(data, passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt config: %v (check your master password)", err)
		}
		data = decrypted
	} else if passphrase != "" {
		// Data is not encrypted but a passphrase is set
		// Encrypt it immediately to ensure security
		fmt.Printf("\n🔒 Encrypting your configuration file...\n")
		var cfg Config
		err = yaml.Unmarshal(data, &cfg)
		if err == nil {
			err = SaveConfig(&cfg, passphrase)
			if err != nil {
				fmt.Printf("⚠️  Failed to encrypt config: %v\n", err)
			} else {
				fmt.Printf("✅ Config encrypted successfully.\n")
			}
		}
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Connections == nil {
		cfg.Connections = make(map[string]Connection)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config, passphrase string) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if passphrase != "" {
		encrypted, err := utils.Encrypt(data, passphrase)
		if err != nil {
			return fmt.Errorf("failed to encrypt config: %v", err)
		}
		data = encrypted
	}

	return os.WriteFile(path, data, 0600)
}

func (cfg *Config) UpdateLastUsed(name string) {
	if conn, ok := cfg.Connections[name]; ok {
		conn.LastUsed = time.Now()
		conn.UsageCount++
		cfg.Connections[name] = conn
	}
}

func (cfg *Config) DeleteConnection(name string) bool {
	if _, ok := cfg.Connections[name]; ok {
		delete(cfg.Connections, name)
		return true
	}
	return false
}

func (cfg *Config) ToggleFavorite(name string) bool {
	if conn, ok := cfg.Connections[name]; ok {
		conn.Favorite = !conn.Favorite
		cfg.Connections[name] = conn
		return conn.Favorite
	}
	return false
}

func (cfg *Config) SetNotes(name, notes string) bool {
	if conn, ok := cfg.Connections[name]; ok {
		conn.Notes = notes
		cfg.Connections[name] = conn
		return true
	}
	return false
}
