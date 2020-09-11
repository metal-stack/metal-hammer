package cmd

import (
	"context"
	"github.com/metal-stack/go-hal/pkg/api"
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

// FetchSupermetalPassword tries to fetch the bmc superuser password from metla-api.
// If no superuser password has been set in metal-api it returns an empty string and true as
// the second return value, which indicates to skip further processing regarding the superuser password.
// Otherwise that second return value is always false.
func (c *GrpcClient) FetchSupermetalPassword() (string, bool, error) {
	client, closer, err := c.newSupermetalPasswordClient()
	if err != nil {
		return "", false, err
	}
	defer closer.Close()

	req := &v1.SupermetalPasswordRequest{}
	resp, err := client.FetchSupermetalPassword(context.Background(), req)
	if err != nil {
		return "", false, err
	}

	if resp.GetFeatureDisabled() {
		return "", true, nil
	}

	return resp.GetSupermetalPassword(), false, nil
}

func (h *Hammer) CreateBmcSuperuser() (bool, error) {
	pwd, skip, err := h.GrpcClient.FetchSupermetalPassword()
	if err != nil {
		return false, errors.Wrap(err, "failed to fetch supermetal password")
	}

	if skip {
		return false, nil
	}

	err = h.Hal.BMCCreateUser(h.Hal.BMCSuperUser(), api.AdministratorPrivilege, pwd)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create bmc superuser: %s", h.Hal.BMCSuperUser().Name)
	}

	return true, nil
}
