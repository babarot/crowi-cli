package cli

import (
	"context"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/b4b4r07/crowi/config"
	"github.com/b4b4r07/crowi/util"
	"github.com/crowi/go-crowi"
)

type Screen struct {
	Text string
}

func NewScreen() (*Screen, error) {
	client, err := GetClient()
	if err != nil {
		return &Screen{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := config.Conf.Crowi.User
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

	text := ""
	for _, page := range res.Pages {
		text += fmt.Sprintf("%s\n", page.Path)
	}

	return &Screen{
		Text: text,
	}, nil
}

func (s *Screen) Filter() (selectedLines []string, err error) {
	lines, err := util.Filter(s.Text)
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

type Crowi struct {
	Path, URL string
}

func ParseLine(line string) (*Crowi, error) {
	return &Crowi{
		Path: line,
		URL:  path.Join(config.Conf.Crowi.BaseURL, line),
	}, nil
}
