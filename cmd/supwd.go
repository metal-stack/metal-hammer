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

	changeIsNeeded, err := bmcConn.NeedsPasswordChange(bmcConn.SuperUser(), resp.SuperUserPassword)
	if err != nil {
		if !changeIsNeeded {
			// metal-hammer should fail
			// it is not the case, that the user is missing or the password is already correct for the superuser
			h.log.Error("failed to verify password change for bmc superuser", "user", bmcConn.SuperUser().Name, "error", err)
		}

		// Proceed: Password need to be updated or superuser created initially
		if err := bmcConn.CreateUser(bmcConn.SuperUser(), api.AdministratorPrivilege, resp.SuperUserPassword); err != nil {
			stillNeeded, e := bmcConn.NeedsPasswordChange(bmcConn.SuperUser(), resp.SuperUserPassword)
			if e != nil || stillNeeded {
				return fmt.Errorf("failed to configure bmc superuser %s: %w", bmcConn.SuperUser().Name, e)
			}
			h.log.Warn("bmc reported an error configuring superuser, but password verification succeeded, proceed",
				"user", bmcConn.SuperUser().Name, "error", err)
		} else {
			h.log.Info("created superuser", "user", bmcConn.SuperUser().Name)
		}
		return nil
	}

	h.log.Info("superuser already present and password correct", "user", bmcConn.SuperUser().Name)
	return nil
}
