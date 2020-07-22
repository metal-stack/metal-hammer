package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	log "github.com/inconshreveable/log15"
	v1 "github.com/metal-stack/metal-hammer/cmd/api/v1"
	"github.com/metal-stack/metal-hammer/metal-core/client/certs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io"
	"time"
)

// Wait until a machine create request was fired
func (h *Hammer) WaitForInstallation(uuid string) error {
	params := certs.NewGrpcClientCertParams()
	resp, err := h.CertsClient.GrpcClientCert(params)
	if err != nil {
		return err
	}

	clientCert, err := tls.X509KeyPair([]byte(resp.Payload.Cert), []byte(resp.Payload.Key))
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(resp.Payload.CaCert))
	if !ok {
		return errors.New("bad certificate")
	}
	c := NewClient(resp.Payload.Address, clientCert, caCertPool)
	defer c.Close()
	c.WaitForInstallation(uuid)
	return nil
}

type Client struct {
	v1.WaitClient
	conn *grpc.ClientConn
}

func NewClient(addr string, clientCert tls.Certificate, caCertPool *x509.CertPool) *Client {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	}
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Error("can not connect with server", "error", err)
	}

	c := &Client{
		WaitClient: v1.NewWaitClient(conn),
		conn:       conn,
	}

	return c
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) WaitForInstallation(machineID string) {
	req := &v1.WaitRequest{
		MachineID: machineID,
	}

	for {
		stream, err := c.Wait(context.Background(), req)
		if err != nil {
			log.Error("failed waiting for installation, retry in 2sec", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for {
			_, err := stream.Recv()
			if err == io.EOF {
				log.Info("machine has been requested for installation", "machineID", machineID)
				return
			}

			if err != nil {
				log.Error("failed waiting for installation, retry in 2sec", "error", err)
				time.Sleep(2 * time.Second)
				break
			}

			log.Info("wait for installation...", "machineID", machineID)
		}
	}
}
