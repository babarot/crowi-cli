package cli

import (
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
