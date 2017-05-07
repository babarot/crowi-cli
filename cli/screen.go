package cli

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/b4b4r07/crowi/api" // TODO (spinner)
	"github.com/crowi/go-crowi"
)

type Screen struct {
	Text  string
	ID    func(string) string
	Pages *crowi.Pages
}

func NewScreen() (*Screen, error) {
	s := api.NewSpinner("Fetching")
	s.Start()
	defer s.Stop()

	client, err := GetClient()
	if err != nil {
		return &Screen{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := Conf.Crowi.User
	if user == "" {
		return &Screen{}, errors.New("user is not defined")
	}

	res, err := client.Pages.List(ctx, "", user, &crowi.PagesListOptions{
		crowi.ListOptions{Pagenation: true},
	})
	if err != nil {
		return &Screen{}, err
	}

	if !res.OK {
		return &Screen{}, errors.New(res.Error)
	}

	// TODO: map ?
	text := ""
	for _, pi := range res.Pages {
		text += fmt.Sprintf("%s\n", pi.Path)
	}

	id := func(path string) (id string) {
		for _, pi := range res.Pages {
			if pi.Path == path {
				return pi.ID
			}
		}
		return ""
	}

	return &Screen{
		Text:  text,
		ID:    id,
		Pages: res,
	}, nil
}

func (s *Screen) Filter() (selectedLines []string, err error) {
	lines, err := Filter(Conf.Core.SelectCmd, s.Text)
	if err != nil {
		return
	}
	for _, line := range lines {
		if line == "" {
			continue
		}
		selectedLines = append(selectedLines, line)
	}
	if len(selectedLines) == 0 {
		return []string{}, errors.New("no lines selected")
	}
	return
}

type Line struct {
	Path      string
	URL       string
	ID        string
	LocalPath string
}

func (s *Screen) ParseLine(line string) (*Line, error) {
	return &Line{
		Path:      line,
		URL:       path.Join(Conf.Crowi.BaseURL, line),
		ID:        s.ID(line),
		LocalPath: filepath.Join(Conf.Crowi.LocalPath, s.ID(line)),
	}, nil
}
