package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/crowi/go-crowi"
)

type Page struct {
	Info      crowi.PageInfo
	LocalPath string
	Client    *crowi.Client
}

func (page Page) Create(path, body string) (*crowi.Page, error) {
	s := NewSpinner("Posting...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return page.Client.Pages.Create(ctx, path, body)
}

func fileContent(fname string) string {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (page Page) upload() (res *crowi.Page, err error) {
	s := NewSpinner("Uploading...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err = page.Client.Pages.Get(ctx, page.Info.Path)
	if err != nil {
		return
	}

	remoteBody := res.Page.Revision.Body
	localBody := fileContent(page.LocalPath)

	if remoteBody == localBody {
		// do nothing
		return &crowi.Page{}, nil
	}

	res, err = page.Client.Pages.Update(ctx, page.Info.ID, localBody)
	return
}

func (page Page) download() (res *crowi.Page, err error) {
	s := NewSpinner("Downloading...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err = page.Client.Pages.Get(ctx, page.Info.Path)
	if err != nil {
		return
	}

	remoteBody := res.Page.Revision.Body
	localBody := fileContent(page.LocalPath)

	if remoteBody == localBody {
		// do nothing
		return &crowi.Page{}, nil
	}

	err = ioutil.WriteFile(page.LocalPath, []byte(remoteBody), os.ModePerm)
	return
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (page Page) Sync() (err error) {
	var res *crowi.Page

	if !exists(page.LocalPath) {
		err = ioutil.WriteFile(page.LocalPath, []byte(page.Info.Revision.Body), os.ModePerm)
		if err != nil {
			return err
		}
	}
	fi, err := os.Stat(page.LocalPath)
	if err != nil {
		return err
	}

	local := fi.ModTime().UTC()
	remote := page.Info.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		res, err = page.upload()
		if res.OK {
			fmt.Printf("Uploaded %s\n", res.Page.Path)
		}
	case remote.After(local):
		res, err = page.download()
		if res.OK {
			fmt.Printf("Downloaded %s\n", res.Page.Path)
		}
	default:
	}

	return err
}
