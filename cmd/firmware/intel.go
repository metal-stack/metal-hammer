package firmware

import (
	"fmt"

	log "github.com/inconshreveable/log15"
)

type intel struct {
	name           string
	desiredVersion string
}

func (r intel) String() string {
	return r.name
}

// firmware update via
// /intel/nvmupdate64e -u -s
func (r intel) update() error {
	output, err := run("/intel/nvmupdate64e", "-u", "-s", "-a", "/intel")
	if err != nil {
		return fmt.Errorf("unable to update intel firmware %w", err)
	}
	log.Info("intel", "updated firware output", output)
	return nil
}

func (r intel) current() (string, error) {
	log.Info("not implemented")
	return "", nil
}

func (r intel) desired() string {
	return r.desiredVersion
}

func (r intel) updateRequired() bool {
	log.Info("not implemented")
	return true
}
