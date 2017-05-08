package cmd

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/b4b4r07/crowi/api"
	"github.com/b4b4r07/crowi/cli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [FILE/DIR]",
	Short: "Create a new page",
	Long:  `Create a new page`,
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

	page := api.Page{Client: client}
	res, err := page.Create(p.path, p.body)
	if err != nil {
		return err
	}

	if !res.OK {
		return errors.New(res.Error)
	}

	cli.Underline("Created", res.Page.ID)
	return nil
}

func makeFromEditor() (p page, err error) {
	user := cli.Conf.Crowi.User
	if user == "" {
		return p, errors.New("config user not defined")
	}
	cli.ScanDefaultString = path.Join(
		"/user", user, "memo", time.Now().Format("2006/01/02"),
	) + "/"

	path, err := cli.Scan(color.YellowString("Path> "), !cli.ScanAllowEmpty)
	if err != nil {
		return
	}
	if !filepath.HasPrefix(path, "/") {
		return page{}, errors.New("path: it must start with a slash")
	}
	// Do not make it a portal page
	path = strings.TrimSuffix(path, "/")

	f, err := cli.TempFile(filepath.Base(path) + cli.Extention)
	defer os.Remove(f.Name())

	editor := cli.Conf.Core.Editor
	if editor == "" {
		return p, errors.New("config editor not defined")
	}
	err = cli.Run(editor, f.Name())
	if err != nil {
		return
	}

	return page{
		path: path,
		body: cli.FileContent(f.Name()),
	}, nil
}

func init() {
	RootCmd.AddCommand(newCmd)
}
