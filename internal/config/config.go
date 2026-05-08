package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/atomic"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB DbConfig `yaml:"DB"`
}

type DbConfig struct {
	DbURL               string        `yaml:"DbURL"`
	Query_timeout_read  time.Duration `yaml:"query_timeout_read"`
	Query_timeout_write time.Duration `yaml:"query_timeout_write"`
	ConnectionInterval  time.Duration `yaml:"connectionInterval"`
	MaxRetry            int           `yaml:"maxRetry"`
	MaxDelay            time.Duration `yaml:"maxDelay"`
}

// Load loads the values frrom the file "path" to the struct c, if the file is not present:
// the default values are loaded and the file is created.
func (c *Config) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.SetDefault()
			c.Save(path)
			return nil
		}
		return fmt.Errorf("Error while reading the config file: %w", err)
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Error during parsing of YAML config file: %w", err)
	}
	return nil
}

// SetDefault sets the config default values
func (c *Config) SetDefault() {
	c.DB.DbURL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"
	c.DB.Query_timeout_read = 5 * time.Second
	c.DB.Query_timeout_write = 2 * time.Second
	c.DB.ConnectionInterval = 3 * time.Second
	c.DB.MaxRetry = 0
	c.DB.MaxDelay = 1 * time.Minute
}

// Save saves the configs to the "path" using the WriteFileAtomic function
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error while parsing to YAML: %w", err)
	}
	return WriteFileAtomic(path, data, 0644)
}

// Data validation after loading
func (c *Config) Validate() error {
	if c.DB.DbURL == "" {
		return fmt.Errorf("[config][db] connection url cannot be empty")
	}
	return nil
}

// WriteFileAtomic saves the file "path" using an atomic sequence.
// First it creates the path if not existing it creates it then calls the write file atomic.
// first it saves the "data" to a tempFile, then it updates the permissions
// of the tempFile to "perm" and then it replaces the original file "path" with the tempFile. This is usefull for files that are core to the app
// and can't risk to be corrupted.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	//I get the directory from the path
	dir := filepath.Dir(path)

	//create the directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Cannot create directory: %w", err)
	}

	// create an io.reader for atomic.WriteFile
	reader := bytes.NewReader(data)

	//atomicaly writing the file
	if err := atomic.WriteFile(path, reader); err != nil {
		return fmt.Errorf("Atomic write failed: %w", err)
	}

	//set permissions
	if err := os.Chmod(path, perm); err != nil {
		return fmt.Errorf("Error setting permissions: %w", err)
	}

	return nil
}
