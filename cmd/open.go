package cmd

import (
	"github.com/b4b4r07/crowi/cli"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a crowi page",
	Long:  "Open a crowi page",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}

	lines, err := screen.Select()
	if err != nil {
		return err
	}

	// TODO: lines (range)
	return cli.OpenURL(lines[0].URL)
}

func init() {
	RootCmd.AddCommand(openCmd)
}
