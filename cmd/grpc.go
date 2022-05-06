package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/metal-core/client/certs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	log  *zap.SugaredLogger
	conn grpc.ClientConnInterface
}

// NewGrpcClient fetches the address and certificates from metal-core needed to communicate with metal-api via grpc,
// and returns a new grpc client that can be used to invoke all provided grpc endpoints.
func NewGrpcClient(log *zap.SugaredLogger, certsClient certs.ClientService) (*GrpcClient, error) {
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

	grpcOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, resp.Payload.Address, grpcOpts...)
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		log:  log,
		conn: conn,
	}, nil
}

func (c *GrpcClient) Event() v1.EventServiceClient {
	return v1.NewEventServiceClient(c.conn)
}

func (c *GrpcClient) Wait() v1.WaitClient {
	return v1.NewWaitClient(c.conn)
}

func (c *GrpcClient) SuperUserPassword() v1.SuperUserPasswordClient {
	return v1.NewSuperUserPasswordClient(c.conn)
}
