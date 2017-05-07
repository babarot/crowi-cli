package api

import (
	"context"
	"fmt"
	"time"

	"github.com/crowi/go-crowi"
)

func CreatePage(cli *crowi.Client, path, body string) (*crowi.Page, error) {
	s := NewSpinner("Posting")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return cli.Pages.Create(ctx, path, body)
}

// func upload() error
// func download() error

type PageData struct {
	Page      crowi.PageInfo
	LocalPath string
}

func SyncPage(data PageData) error {
	fmt.Printf("Synced %#v\n", data.Page)
	return nil
}
