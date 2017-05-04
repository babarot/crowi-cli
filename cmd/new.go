package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/b4b4r07/crowi/config"
	"github.com/b4b4r07/crowi/util"
	"github.com/crowi/go-crowi"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [FILE/DIR]",
	Short: "Create a new gist",
	Long:  `Create a new gist. If you pass file/dir paths, upload those files`,
	RunE:  new,
}

type gistItem struct {
	path, body string
}

func new(cmd *cobra.Command, args []string) error {
	var err error
	var gi gistItem

	gi, err = makeFromEditor()
	if err != nil {
		return err
	}

	client, err := crowi.NewClient(config.Conf.Core.BaseURL, config.Conf.Core.Token)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.PagesCreate(ctx, gi.path, gi.body)
	if err != nil {
		return err
	}

	util.Underline("Created", res.Page.ID)
	return nil
}

func makeFromEditor() (gi gistItem, err error) {
	util.ScanDefaultString = fmt.Sprintf("/user/%s/memo/%s/", os.Getenv("USER"), time.Now().Format("2006/01/02"))
	filename, err := util.Scan(color.YellowString("Filename> "), !util.ScanAllowEmpty)
	if err != nil {
		return
	}
	f, err := util.TempFile(filepath.Base(filename))
	defer os.Remove(f.Name())
	err = util.RunCommand(config.Conf.Core.Editor, f.Name())
	if err != nil {
		return
	}
	return gistItem{
		path: filename,
		body: util.FileContent(f.Name()),
	}, nil
}

func init() {
	RootCmd.AddCommand(newCmd)
}
