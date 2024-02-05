package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log/slog"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	metalgo "github.com/metal-stack/metal-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type MetalAPIClient struct {
	log    *slog.Logger
	conn   grpc.ClientConnInterface
	Driver metalgo.Client
}

// NewMetalAPIClient fetches the address,hmac and certificates from pixie needed to communicate with metal-api,
// and returns a new client that can be used to invoke all provided grpc and rest endpoints.
func NewMetalAPIClient(log *slog.Logger, spec *Specification) (*MetalAPIClient, error) {
	metalConfig := spec.MetalConfig

	clientCert, err := tls.X509KeyPair([]byte(metalConfig.Cert), []byte(metalConfig.Key))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(metalConfig.CACert))
	if !ok {
		return nil, errors.New("bad certificate")
	}

	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
		MinVersion:   tls.VersionTLS12,
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, metalConfig.GRPCAddress, grpcOpts...)
	if err != nil {
		return nil, err
	}

	driver, err := metalgo.NewDriver(metalConfig.MetalAPIUrl, "", metalConfig.HMAC, metalgo.AuthType("Metal-View"))
	if err != nil {
		return nil, err
	}

	return &MetalAPIClient{
		log:    log,
		conn:   conn,
		Driver: driver,
	}, nil
}

func (c *MetalAPIClient) Event() v1.EventServiceClient {
	return v1.NewEventServiceClient(c.conn)
}

func (c *MetalAPIClient) BootService() v1.BootServiceClient {
	return v1.NewBootServiceClient(c.conn)
}
