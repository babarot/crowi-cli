package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core Core
}

type Core struct {
	Token     string `toml:"token"`
	BaseURL   string `toml:"base_url"`
	Editor    string `toml:"editor"`
	TomlFile  string `toml:"toml_file"`
	SelectCmd string `toml:"selectcmd"`
}

var Conf Config

func GetDefaultDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}
	dir = filepath.Join(dir, "crowi")

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return dir, fmt.Errorf("cannot create directory: %v", err)
	}

	return dir, nil
}

func (cfg *Config) LoadFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.Core.Token = os.Getenv("CROWI_ACCESS_TOKEN")
	cfg.Core.BaseURL = "https://wiki.your.domain"
	cfg.Core.Editor = "vim"
	cfg.Core.TomlFile = file
	cfg.Core.SelectCmd = "fzf"

	return toml.NewEncoder(f).Encode(cfg)
}
