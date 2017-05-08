package cmd

import (
	"github.com/b4b4r07/crowi/cli"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a page",
	Long:  `Edit a page`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}
	lines, err := screen.Select()
	if err != nil {
		return err
	}

	// TODO: lines (range)
	return cli.EditPage(screen.Pages, lines[0])
}

func init() {
	RootCmd.AddCommand(editCmd)
}
