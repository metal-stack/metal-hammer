package cmd

import (
	"context"
	"fmt"

	"github.com/metal-stack/go-hal/pkg/api"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
)

// createBmcSuperuser creates the bmc super user.
func (h *hammer) createBmcSuperuser() error {
	req := &v1.BootServiceSuperUserPasswordRequest{}
	resp, err := h.metalAPIClient.BootService().SuperUserPassword(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to fetch SuperUser password %w", err)
	}

	if resp.SuperUserPassword == "" {
		h.log.Warn("creation of superuser disabled because password is empty")
		return nil
	}

	bmcConn := h.hal.BMCConnection()

	err = bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, resp.SuperUserPassword)
	if err != nil {
		// FIXME: this happens always after the first creation on X12 and newer boards
		// return fmt.Errorf("failed to create bmc superuser: %s %w", bmcConn.SuperUser().Name, err)
		h.log.Error("failed to create bmc superuser", "user", bmcConn.SuperUser().Name, "error", err)
		return nil
	}

	h.log.Info("created superuser", "user", bmcConn.SuperUser().Name)
	return nil
}
