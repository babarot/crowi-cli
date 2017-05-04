package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/b4b4r07/crowi/config"
	"github.com/crowi/go-crowi"
)

type Screen struct {
	Text string
}

func NewScreen() (*Screen, error) {
	client, err := crowi.NewClient(config.Conf.Core.BaseURL, config.Conf.Gist.Token)
	if err != nil {
		return &Screen{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.PagesList(ctx, "", "b4b4r07")
	if err != nil {
		return &Screen{}, err
	}

	text := ""
	for _, page := range res.Pages {
		text += fmt.Sprintf("%s\n", page.Path)
	}

	return &Screen{
		Text: text,
	}, nil
}
