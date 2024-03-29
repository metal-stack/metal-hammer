package firmware

import "log/slog"

type raidcontroller struct {
	name           string
	desiredVersion string
	log            *slog.Logger
}

func (r raidcontroller) String() string {
	return r.name
}

// firmware update via
// storcli /cX download file=smc3108.rom
func (r raidcontroller) update() error {
	r.log.Error("not implemented")
	return nil
}

func (r raidcontroller) current() (string, error) {
	r.log.Error("not implemented")
	return "", nil
}

func (r raidcontroller) desired() string {
	return r.desiredVersion
}

func (r raidcontroller) updateRequired() bool {
	r.log.Error("not implemented")
	return true
}
