package cmd

import (
	"context"
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/go-hal/pkg/api"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/pkg/errors"
	"io"
)

func (c *GrpcClient) newSuperUserPasswordClient() (v1.SuperUserPasswordClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewSuperUserPasswordClient(conn), conn, nil
}

// FetchSuperUserPassword tries to fetch the bmc superuser password from metla-api.
// If no superuser password has been set in metal-api it returns an empty string and true as
// the second return value, which indicates to skip further processing regarding the superuser password.
// Otherwise that second return value is always false.
func (c *GrpcClient) FetchSuperUserPassword() (string, bool, error) {
	client, closer, err := c.newSuperUserPasswordClient()
	if err != nil {
		return "", false, err
	}
	defer closer.Close()

	req := &v1.SuperUserPasswordRequest{}
	resp, err := client.FetchSuperUserPassword(context.Background(), req)
	if err != nil {
		return "", false, err
	}

	if resp.GetFeatureDisabled() {
		return "", true, nil
	}

	return resp.GetSuperUserPassword(), false, nil
}

func (h *Hammer) CreateBmcSuperuser() error {
	pwd, featureDisabled, err := h.GrpcClient.FetchSuperUserPassword()
	if err != nil {
		return errors.Wrap(err, "failed to fetch SuperUser password")
	}

	if featureDisabled {
		return nil
	}

	bmcConn := h.Hal.BMCConnection()
	err = bmcConn.SetUserEnabled(bmcConn.PresentSuperUser(), false)
	if err != nil {
		log.Error("failed to disable present bmc admin user", "error", err)
		return err
	}

	err = bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, pwd)
	if err != nil {
		return errors.Wrapf(err, "failed to create bmc superuser: %s", bmcConn.SuperUser().Name)
	}

	return nil
}
