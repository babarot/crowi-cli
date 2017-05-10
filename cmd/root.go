package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/b4b4r07/crowi/cli"
	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var showVersion bool

var RootCmd = &cobra.Command{
	Use:           "crowi",
	Short:         "crowi command-line interface",
	Long:          "crowi - A simple Crowi editor for CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          noArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("version %s/%s\n", Version, runtime.Version())
			return
		}
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	initConf()
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show the version and exit")
}

func initConf() {
	dir, _ := cli.GetDefaultDir()
	toml := filepath.Join(dir, "config.toml")

	err := cli.Conf.LoadFile(toml)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	suggestionsString := ""
	if !cmd.DisableSuggestions {
		if cmd.SuggestionsMinimumDistance <= 0 {
			cmd.SuggestionsMinimumDistance = 2
		}
		suggestions := cmd.SuggestionsFor(args[0])
		switch len(suggestions) {
		case 0:
			// Ignore because crowi-XXX may be present
			break
		// case 1:
		// TODO: call function dynamically by name
		// 	fmt.Println(suggestions[0])
		default:
			suggestionsString += "\n\nDid you mean this?\n"
			for _, s := range suggestions {
				suggestionsString += fmt.Sprintf("\t%v\n", s)
			}
			return fmt.Errorf("unknown command %q for crowi%s", args[0], suggestionsString)
		}
	}
	crowiCmd := cmd.CommandPath() + "-" + args[0]
	if _, err := exec.LookPath(crowiCmd); err != nil {
		return fmt.Errorf(
			"crowi: '%s' is not a crowi command.\nSee 'crowi --help'", args[0])
	}
	out, err := exec.Command(crowiCmd).Output()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		fmt.Printf(string(out))
	}
	return nil
}
