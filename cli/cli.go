package cli

import (
	"github.com/b4b4r07/crowi/config"
	"github.com/crowi/go-crowi"
)

func GetClient() (*crowi.Client, error) {
	return crowi.NewClient(
		crowi.Config{
			URL:   config.Conf.Crowi.BaseURL,
			Token: config.Conf.Crowi.Token,
		})
}
