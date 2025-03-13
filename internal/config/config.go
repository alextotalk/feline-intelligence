package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Env     string `yaml:"env"`
	} `yaml:"app"`

	Server struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		ReadTimeout  string `yaml:"read_timeout"`
		WriteTimeout string `yaml:"write_timeout"`
		IdleTimeout  string `yaml:"idle_timeout"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		log.Fatal("Config path is empty")
		return nil, os.ErrInvalid
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
		return nil, err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", absPath)
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(absPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
