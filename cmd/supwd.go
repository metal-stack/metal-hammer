package cmd

import (
	"context"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/pkg/errors"
	"io"
)

func (c *GrpcClient) newSupermetalPasswordClient() (v1.SupermetalPasswordClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewSupermetalPasswordClient(conn), conn, nil
}

func (c *GrpcClient) FetchSupermetalPassword(partitionID string) (string, error) {
	client, closer, err := c.newSupermetalPasswordClient()
	if err != nil {
		return "", err
	}
	defer closer.Close()

	req := &v1.SupermetalPasswordRequest{
		PartitionID: partitionID,
	}
	resp, err := client.FetchSupermetalPassword(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.GetSupermetalPassword(), nil
}

func (h *Hammer) UpdateBmcSuperuserPassword(partitionID string) error {
	supwd, err := h.GrpcClient.FetchSupermetalPassword(partitionID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch supermetal password")
	}

	err = h.Hal.BMCChangePassword(h.Hal.BMCSuperUser(), supwd)
	if err != nil {
		return errors.Wrap(err, "failed to change bmc superuser password")
	}

	return nil
}
