package firmware

import (
	"context"
	"log/slog"
)

// raidcontroller is not implemented yet and therefore not registered in New,
// the assertion keeps it in sync with the updater interface.
var _ updater = raidcontroller{}

type raidcontroller struct {
	name           string
	desiredVersion string
	log            *slog.Logger
}

func (r raidcontroller) String() string {
	return r.name
}

func (r raidcontroller) applicable(ctx context.Context) bool {
	r.log.Error("not implemented")
	return false
}

// firmware update via
// storcli /cX download file=smc3108.rom
// Unlike the bank switching nvm of the intel cards, a raid controller runs the firmware
// which was written to it after the next reboot, so once implemented this returns
// ActivationWarmReboot when it flashed a controller.
func (r raidcontroller) update(ctx context.Context) (Activation, error) {
	r.log.Error("not implemented")
	return ActivationNone, nil
}

func (r raidcontroller) current(ctx context.Context) (string, error) {
	r.log.Error("not implemented")
	return "", nil
}

func (r raidcontroller) desired(ctx context.Context) (string, error) {
	return r.desiredVersion, nil
}
