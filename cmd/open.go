package cmd

import (
	"errors"

	"github.com/b4b4r07/crowi/cli"
	"github.com/b4b4r07/crowi/util"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open user's gist",
	Long:  "Open user's gist",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
	s, err := cli.NewScreen()
	if err != nil {
		return err
	}

	selectedLines, err := util.Filter(s.Text)
	if err != nil {
		return err
	}

	if len(selectedLines) == 0 {
		return errors.New("no page selected")
	}

	line, err := cli.ParseLine(selectedLines[0])
	if err != nil {
		return err
	}

	return util.Open(line.URL)
}

func init() {
	RootCmd.AddCommand(openCmd)
}
