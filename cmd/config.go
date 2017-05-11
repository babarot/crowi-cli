package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/b4b4r07/crowi/cli"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "config",
	Short: "Config the setting file",
	Long:  "Config the setting file with your editor (default: vim)",
	RunE:  conf,
}

func conf(cmd *cobra.Command, args []string) error {
	if confGetKey != "" {
		dir, err := cli.GetDefaultDir()
		if err != nil {
			return err
		}
		config, err := toml.LoadFile(filepath.Join(dir, "config.toml"))
		if err != nil {
			return err
		}
		value := config.Get(confGetKey)
		if value != nil {
			fmt.Printf("%v\n", value)
			return nil
		}
		return fmt.Errorf("%s: no such key found", confGetKey)
	}

	editor := cli.Conf.Core.Editor
	tomlfile := cli.Conf.Core.TomlFile
	if tomlfile == "" {
		dir, _ := cli.GetDefaultDir()
		tomlfile = filepath.Join(dir, "config.toml")
	}
	return cli.Run(editor, tomlfile)
}

var confGetKey string

func init() {
	RootCmd.AddCommand(confCmd)
	confCmd.Flags().StringVarP(&confGetKey, "get", "", "", "Get the config value")
}
