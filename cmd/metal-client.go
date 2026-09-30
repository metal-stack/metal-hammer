package cmd

import (
	"log/slog"

	"github.com/metal-stack/api/go/client"
	pixiecore "github.com/metal-stack/pixie/api"
)

func NewMetalAPIClient(log *slog.Logger, clientConfig *pixiecore.Client) (client.Client, error) {
	client, err := client.New(&client.DialConfig{
		BaseURL: clientConfig.ApiUrl,
		Token:   clientConfig.Token,
		Log:     log,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
