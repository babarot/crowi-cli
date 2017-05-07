package api

import (
	"context"
	"time"

	"github.com/b4b4r07/crowi/cli"
	"github.com/crowi/go-crowi"
)

func PagesCreate(path, body string) (*crowi.Page, error) {
	s := cli.NewSpinner("Posting")
	s.Start()
	defer s.Stop()

	client, err := cli.GetClient()
	if err != nil {
		return &crowi.Page{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return client.Pages.Create(ctx, path, body)
}
