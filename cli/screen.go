package cli

import (
	"context"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/b4b4r07/crowi/util"
	"github.com/briandowns/spinner"
	"github.com/crowi/go-crowi"
)

type Screen struct {
	Text string
}

var (
	SpinnerSymbol int = 14
)

func NewScreen() (*Screen, error) {
	s := spinner.New(spinner.CharSets[SpinnerSymbol], 100*time.Millisecond)
	s.Suffix = " Fetching..."
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

	text := ""
	for _, page := range res.Pages {
		text += fmt.Sprintf("%s\n", page.Path)
	}

	return &Screen{
		Text: text,
	}, nil
}

func (s *Screen) Filter() (selectedLines []string, err error) {
	lines, err := util.Filter(Conf.Core.SelectCmd, s.Text)
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

type Page struct {
	Path, URL string
}

func ParseLine(line string) (*Page, error) {
	return &Page{
		Path: line,
		URL:  path.Join(Conf.Crowi.BaseURL, line),
	}, nil
}
