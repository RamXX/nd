package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RamXX/vlt"
	"gopkg.in/yaml.v3"
)

// Config holds vault-level nd configuration stored in .nd.yaml.
type Config struct {
	Version   string `yaml:"version"`
	Prefix    string `yaml:"prefix"`
	CreatedBy string `yaml:"created_by"`
}

// Store wraps a vlt.Vault with issue-tracker operations.
type Store struct {
	vault  *vlt.Vault
	config Config
	dir    string
}

// Open opens an existing nd vault at dir.
func Open(dir string) (*Store, error) {
	v, err := vlt.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open vault: %w", err)
	}
	s := &Store{vault: v, dir: dir}
	if err := s.loadConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return s, nil
}

// Init creates a new nd vault at dir.
func Init(dir, prefix, author string) (*Store, error) {
	// Create vault directory structure.
	for _, sub := range []string{"issues", ".trash"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", sub, err)
		}
	}

	// Write config.
	cfg := Config{
		Version:   "1",
		Prefix:    prefix,
		CreatedBy: author,
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".nd.yaml"), data, 0o644); err != nil {
		return nil, fmt.Errorf("write config: %w", err)
	}

	v, err := vlt.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open vault after init: %w", err)
	}
	return &Store{vault: v, config: cfg, dir: dir}, nil
}

func (s *Store) loadConfig() error {
	data, err := os.ReadFile(filepath.Join(s.dir, ".nd.yaml"))
	if err != nil {
		return fmt.Errorf("read .nd.yaml: %w", err)
	}
	return yaml.Unmarshal(data, &s.config)
}

// Vault returns the underlying vlt.Vault for direct operations.
func (s *Store) Vault() *vlt.Vault { return s.vault }

// Config returns the nd configuration.
func (s *Store) Config() Config { return s.config }

// Dir returns the vault root directory.
func (s *Store) Dir() string { return s.dir }

// Prefix returns the configured issue ID prefix.
func (s *Store) Prefix() string { return s.config.Prefix }

// IssueExists checks whether an issue with the given ID exists in the vault.
func (s *Store) IssueExists(id string) bool {
	p := filepath.Join(s.dir, "issues", id+".md")
	_, err := os.Stat(p)
	return err == nil
}
