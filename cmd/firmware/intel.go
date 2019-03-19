package firmware

import (
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
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
	output, err := run("nvmupdate64e", "-u", "-s")
	if err != nil {
		return errors.Wrap(err, "unable to update intel firmware")
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
