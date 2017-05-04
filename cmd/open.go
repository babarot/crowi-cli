package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/b4b4r07/crowi/config"
	"github.com/b4b4r07/crowi/util"
	"github.com/crowi/go-crowi"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open user's gist",
	Long:  "Open user's gist",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
	client, err := crowi.NewClient(config.Conf.Core.BaseURL, config.Conf.Gist.Token)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.PagesList(ctx, "", "b4b4r07")
	if err != nil {
		panic(err)
	}

	text := ""
	for _, page := range res.Pages {
		text += fmt.Sprintf("%s\n", page.Path)
	}

	selectedLines, err := util.Filter(text)
	if err != nil {
		return err
	}

	if len(selectedLines) == 0 {
		return errors.New("No gist selected")
	}

	fmt.Printf("%#v\n", selectedLines)
	return nil
}

func init() {
	RootCmd.AddCommand(openCmd)
}
