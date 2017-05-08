package cli

import (
	"errors"

	"github.com/b4b4r07/crowi/api"
	"github.com/crowi/go-crowi"
)

func GetClient() (*crowi.Client, error) {
	return crowi.NewClient(
		crowi.Config{
			URL:                Conf.Crowi.BaseURL,
			Token:              Conf.Crowi.Token,
			InsecureSkipVerify: true,
		})
}

func EditPage(res *crowi.Pages, line Line) error {
	var (
		err  error
		info crowi.PageInfo
	)

	client, err := GetClient()
	if err != nil {
		return err
	}

	for _, pageInfo := range res.Pages {
		if pageInfo.ID == line.ID {
			info = pageInfo
		}
	}

	page := api.Page{
		Info:      info,
		LocalPath: line.LocalPath,
		Client:    client,
	}

	// sync before editing
	err = page.Sync()
	if err != nil {
		return err
	}

	editor := Conf.Core.Editor
	if editor == "" {
		return errors.New("config editor not set")
	}
	if err = Run(editor, line.LocalPath); err != nil {
		return err
	}

	// sync after editing
	return page.Sync()
}
