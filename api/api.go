package api

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/crowi/go-crowi"
)

type PageData struct {
	Page      crowi.PageInfo
	LocalPath string
	Client    *crowi.Client
}

func (pd PageData) CreatePage(path, body string) (*crowi.Page, error) {
	s := NewSpinner("Posting...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pd.Client.Pages.Create(ctx, path, body)
}

func fileContent(fname string) string {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type (
	uploaded   bool
	downloaded bool
)

func (pd PageData) upload() (done uploaded, err error) {
	s := NewSpinner("Uploading...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := pd.Client.Pages.Get(ctx, pd.Page.Path)
	if err != nil {
		return
	}

	remoteBody := res.Page.Revision.Body
	localBody := fileContent(pd.LocalPath)

	if remoteBody == localBody {
		// do nothing
		return
	}

	_, err = pd.Client.Pages.Update(ctx, pd.Page.ID, localBody)
	return uploaded(true), err
}

func (pd PageData) download() (done downloaded, err error) {
	s := NewSpinner("Downloading...")
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := pd.Client.Pages.Get(ctx, pd.Page.Path)
	if err != nil {
		return
	}

	remoteBody := res.Page.Revision.Body
	localBody := fileContent(pd.LocalPath)

	if remoteBody == localBody {
		// do nothing
		return
	}

	err = ioutil.WriteFile(pd.LocalPath, []byte(remoteBody), os.ModePerm)
	return downloaded(true), err
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (pd PageData) SyncPage() error {
	var (
		done interface{}
		err  error
	)

	if !exists(pd.LocalPath) {
		err := ioutil.WriteFile(pd.LocalPath, []byte(pd.Page.Revision.Body), os.ModePerm)
		if err != nil {
			return err
		}
	}
	fi, err := os.Stat(pd.LocalPath)
	if err != nil {
		return err
	}

	local := fi.ModTime().UTC()
	remote := pd.Page.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		done, err = pd.upload()
	case remote.After(local):
		done, err = pd.download()
	default:
	}

	switch done := done.(type) {
	case uploaded:
		if done {
			println("uploaded")
		}
	case downloaded:
		if done {
			println("downloaded")
		}
	}

	return err
}
