package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/b4b4r07/crowi/cli"
	"github.com/b4b4r07/crowi/config"
	"github.com/b4b4r07/crowi/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [FILE/DIR]",
	Short: "Create a new gist",
	Long:  `Create a new gist. If you pass file/dir paths, upload those files`,
	RunE:  new,
}

type page struct {
	path, body string
}

func new(cmd *cobra.Command, args []string) error {
	p, err := makeFromEditor()
	if err != nil {
		return err
	}

	client, err := cli.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.Pages.Create(ctx, p.path, p.body)
	if err != nil {
		return err
	}

	if !res.OK {
		return errors.New(res.Error)
	}

	util.Underline("Created", res.Page.ID)
	return nil
}

func makeFromEditor() (p page, err error) {
	util.ScanDefaultString = fmt.Sprintf(
		"/user/%s/memo/%s/",
		os.Getenv("USER"), time.Now().Format("2006/01/02"),
	)

	path, err := util.Scan(color.YellowString("Path> "), !util.ScanAllowEmpty)
	if err != nil {
		return
	}
	if !filepath.HasPrefix(path, "/") {
		return page{}, errors.New("invalid format")
	}
	// Do not make it a portal page
	path = strings.TrimSuffix(path, "/")

	f, err := util.TempFile(filepath.Base(path))
	defer os.Remove(f.Name())

	err = util.RunCommand(config.Conf.Core.Editor, f.Name())
	if err != nil {
		return
	}

	return page{
		path: path,
		body: util.FileContent(f.Name()),
	}, nil
}

func init() {
	RootCmd.AddCommand(newCmd)
}
