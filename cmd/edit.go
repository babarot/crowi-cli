package cmd

import (
	"github.com/b4b4r07/crowi/api"
	"github.com/b4b4r07/crowi/cli"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit",
	Long:  `edit`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}
	selectedLines, err := screen.Filter()
	if err != nil {
		return err
	}

	line, err := screen.ParseLine(selectedLines[0])
	if err != nil {
		return err
	}

	return api.EditPage(line.LocalPath)
}

func init() {
	RootCmd.AddCommand(editCmd)
}
