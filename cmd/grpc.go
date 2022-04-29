package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/metal-core/client/certs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	addr     string
	dialOpts []grpc.DialOption
	log      *zap.SugaredLogger
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
	return &GrpcClient{
		addr: resp.Payload.Address,
		log:  log,
		dialOpts: []grpc.DialOption{
			grpc.WithKeepaliveParams(kacp),
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
			grpc.WithBlock(),
		},
	}, nil
}

func (c *GrpcClient) NewEventClient() (v1.EventServiceClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewEventServiceClient(conn), conn, nil
}

func (c *GrpcClient) NewWaitClient() (v1.WaitClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewWaitClient(conn), conn, nil
}

func (c *GrpcClient) newSuperUserPasswordClient() (v1.SuperUserPasswordClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewSuperUserPasswordClient(conn), conn, nil
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
