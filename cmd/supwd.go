package cmd

import (
	"context"
	"fmt"

	"github.com/metal-stack/go-hal/pkg/api"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
)

// createBmcSuperuser creates the bmc super user.
func (h *Hammer) createBmcSuperuser() error {
	req := &v1.BootServiceSuperUserPasswordRequest{}
	resp, err := h.MetalAPIClient.BootService().SuperUserPassword(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to fetch SuperUser password %w", err)
	}

	if resp.GetFeatureDisabled() {
		h.log.Info("creation of superuser disabled")
		return nil
	}

	bmcConn := h.Hal.BMCConnection()

	h.log.Info("create superuser", "user", bmcConn.SuperUser().Name)

	err = bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, resp.SuperUserPassword)
	if err != nil {
		return fmt.Errorf("failed to create bmc superuser: %s %w", bmcConn.SuperUser().Name, err)
	}

	h.log.Info("superuser created")

	return nil
}
