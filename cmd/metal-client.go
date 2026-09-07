package cmd

import (
	"log/slog"

	"github.com/metal-stack/api/go/client"
	"github.com/metal-stack/api/go/metalstack/admin/v2/adminv2connect"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
)

type MetalAPIClient struct {
	log    *slog.Logger
	client client.Client
}

// NewMetalAPIClient fetches the address,hmac and certificates from pixie needed to communicate with metal-api,
// and returns a new client that can be used to invoke all provided grpc and rest endpoints.
func NewMetalAPIClient(log *slog.Logger, spec *Specification) (*MetalAPIClient, error) {
	metalConfig := spec.MetalConfig

	client, err := client.New(&client.DialConfig{
		BaseURL: metalConfig.MetalAPIServerUrl,
		Token:   metalConfig.MetalAPIServerTokenForHammer,
		Log:     log,
	})
	if err != nil {
		return nil, err
	}

	return &MetalAPIClient{
		log:    log,
		client: client,
	}, nil
}
func (c *MetalAPIClient) Machine() adminv2connect.MachineServiceClient {
	return c.client.Adminv2().Machine()
}
func (c *MetalAPIClient) Event() infrav2connect.EventServiceClient {
	return c.client.Infrav2().Event()
}

func (c *MetalAPIClient) BootService() infrav2connect.BootServiceClient {
	return c.client.Infrav2().Boot()
}
