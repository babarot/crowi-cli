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
		page crowi.PageInfo
	)

	for _, pageInfo := range res.Pages {
		if pageInfo.ID == line.ID {
			page = pageInfo
		}
	}

	data := api.PageData{
		Page:      page,
		LocalPath: line.LocalPath,
	}

	// sync before editing
	err = api.SyncPage(data)
	if err != nil {
		return err
	}

	editor := Conf.Core.Editor
	if editor == "" {
		return errors.New("config editor not set")
	}
	err = Run(editor, line.LocalPath)
	if err != nil {
		return err
	}

	// sync after editing
	return api.SyncPage(data)
}
