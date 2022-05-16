package cmd

import (
	"context"
	"fmt"

	"github.com/metal-stack/go-hal/pkg/api"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
)

// createBmcSuperuser creates the bmc super user.
func (h *Hammer) createBmcSuperuser() error {
	req := &v1.SuperUserPasswordRequest{}
	resp, err := h.MetalAPIClient.SuperUserPassword().FetchSuperUserPassword(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to fetch SuperUser password %w", err)
	}

	if resp.FeatureDisabled {
		return nil
	}

	bmcConn := h.Hal.BMCConnection()
	err = bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, resp.SuperUserPassword)
	if err != nil {
		return fmt.Errorf("failed to create bmc superuser: %s %w", bmcConn.SuperUser().Name, err)
	}

	return nil
}
