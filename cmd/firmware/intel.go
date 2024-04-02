package firmware

import (
	"fmt"
	"log/slog"
)

type intel struct {
	name           string
	desiredVersion string
	log            *slog.Logger
}

func (r intel) String() string {
	return r.name
}

// firmware update via
// /intel/nvmupdate64e -u -s
func (r intel) update() error {
	output, err := run(r.log, "/intel/nvmupdate64e", "-u", "-s", "-a", "/intel")
	if err != nil {
		return fmt.Errorf("unable to update intel firmware %w", err)
	}
	r.log.Info("intel", "updated firmware output", output)
	return nil
}

func (r intel) current() (string, error) {
	r.log.Info("not implemented")
	return "", nil
}

func (r intel) desired() string {
	return r.desiredVersion
}

func (r intel) updateRequired() bool {
	r.log.Info("not implemented")
	return true
}
