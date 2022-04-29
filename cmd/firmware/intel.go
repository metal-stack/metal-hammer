package firmware

import (
	"fmt"

	"go.uber.org/zap"
)

type intel struct {
	name           string
	desiredVersion string
	log            *zap.SugaredLogger
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
	r.log.Infow("intel", "updated firware output", output)
	return nil
}

func (r intel) current() (string, error) {
	r.log.Infow("not implemented")
	return "", nil
}

func (r intel) desired() string {
	return r.desiredVersion
}

func (r intel) updateRequired() bool {
	r.log.Infow("not implemented")
	return true
}
