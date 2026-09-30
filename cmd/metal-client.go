package cmd

import (
	"log/slog"

	"github.com/metal-stack/api/go/client"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
	pixiecore "github.com/metal-stack/pixie/api"
)

func NewMetalAPIClient(log *slog.Logger, clientConfig *pixiecore.Client) (infrav2connect.BootServiceClient, error) {
	client, err := client.New(&client.DialConfig{
		BaseURL: clientConfig.ApiUrl,
		Token:   clientConfig.Token,
		Log:     log,
	})
	if err != nil {
		return nil, err
	}

	return client.Infrav2().Boot(), nil
}
