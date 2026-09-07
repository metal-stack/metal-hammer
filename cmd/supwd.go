package cmd

import (
	"context"
	"fmt"

	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/go-hal/pkg/api"
)

// createBmcSuperuser creates the bmc super user.
func (h *hammer) createBmcSuperuser() error {
	resp, err := h.metalAPIClient.BootService().SuperUserPassword(context.Background(), &infrav2.BootServiceSuperUserPasswordRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch SuperUser password %w", err)
	}

	if resp.SuperUserPassword == "" {
		h.log.Warn("creation of superuser disabled because password is empty")
		return nil
	}

	bmcConn := h.hal.BMCConnection()

	changeIsNeeded, err := bmcConn.NeedsPasswordChange(bmcConn.SuperUser(), resp.SuperUserPassword)
	if err != nil {
		if changeIsNeeded {
			err = bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, resp.SuperUserPassword)
			if err != nil {
				// FIXME: this happens always after the first creation on X12 and newer boards
				// return fmt.Errorf("failed to create bmc superuser: %s %w", bmcConn.SuperUser().Name, err)
				h.log.Error("failed to create bmc superuser", "user", bmcConn.SuperUser().Name, "error", err)
				return nil
			}
		}
		h.log.Error("failed to verify password change for bmc superuser", "user", bmcConn.SuperUser().Name, "error", err)
	}

	h.log.Info("created superuser", "user", bmcConn.SuperUser().Name)
	return nil
}
