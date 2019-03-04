package firmware

import (
	log "github.com/inconshreveable/log15"
)

type raidcontroller struct {
	name           string
	desiredVersion string
}

func (r raidcontroller) String() string {
	return r.name
}

// firmware update via
// storcli /cX download file=smc3108.rom
func (r raidcontroller) update() error {
	log.Error("not implemented")
	return nil
}

func (r raidcontroller) current() (string, error) {
	log.Error("not implemented")
	return "", nil
}

func (r raidcontroller) desired() string {
	return r.desiredVersion
}

func (r raidcontroller) updateRequired() bool {
	log.Error("not implemented")
	return true
}
