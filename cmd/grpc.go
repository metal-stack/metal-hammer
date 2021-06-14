package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/metal-core/client/certs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	*event.EventEmitter
	addr     string
	dialOpts []grpc.DialOption
}

// NewGrpcClient fetches the address and certificates from metal-core needed to communicate with metal-api via grpc,
// and returns a new grpc client that can be used to invoke all provided grpc endpoints.
func NewGrpcClient(certsClient certs.ClientService, emitter *event.EventEmitter) (*GrpcClient, error) {
	params := certs.NewGrpcClientCertParams()
	resp, err := certsClient.GrpcClientCert(params)
	if err != nil {
		return nil, err
	}

	clientCert, err := tls.X509KeyPair([]byte(resp.Payload.Cert), []byte(resp.Payload.Key))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(resp.Payload.CaCert))
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
	return &GrpcClient{
		EventEmitter: emitter,
		addr:         resp.Payload.Address,
		dialOpts: []grpc.DialOption{
			grpc.WithKeepaliveParams(kacp),
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
			grpc.WithBlock(),
		},
	}, nil
}

func (c *GrpcClient) newConnection() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.addr, c.dialOpts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}
