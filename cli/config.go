package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core  CoreConfig  `toml:"core"`
	Crowi CrowiConfig `toml:"crowi"`
}

type CoreConfig struct {
	Editor    string `toml:"editor"`
	TomlFile  string `toml:"toml_file"`
	SelectCmd string `toml:"selectcmd"`
}

type CrowiConfig struct {
	Token     string `toml:"token"`
	BaseURL   string `toml:"base_url"`
	User      string `toml:"user"`
	LocalPath string `toml:"local_path"`
	Paging    bool   `toml:"paging"`
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

	dir, err := GetDefaultDir()
	if err != nil {
		return err
	}

	cfg.Core.Editor = os.Getenv("EDITOR")
	cfg.Core.TomlFile = file
	cfg.Core.SelectCmd = "fzf"
	cfg.Crowi.Token = os.Getenv("CROWI_ACCESS_TOKEN")
	cfg.Crowi.BaseURL = "https://wiki.your.domain"
	cfg.Crowi.User = os.Getenv("USER")
	cfg.Crowi.LocalPath = filepath.Join(dir, "pages")

	os.MkdirAll(cfg.Crowi.LocalPath, 0700)

	return toml.NewEncoder(f).Encode(cfg)
}
