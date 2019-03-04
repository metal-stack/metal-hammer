package firmware

import (
	log "github.com/inconshreveable/log15"
)

// updater check if a firmware update is required and updates
// the firmware if required.
type updater interface {
	update() error
	current() (string, error)
	desired() string
	updateRequired() bool
}

// Firmware take care of firmware management
type Firmware struct {
	updaters []updater
}

// New create a new Firmware manager with all Updaters.
func New() *Firmware {

	raid := raidcontroller{
		name:           "lsi3108",
		desiredVersion: "4.680.00-8290",
	}
	return &Firmware{
		updaters: []updater{raid},
	}
}

// Update run updates for all firmwares found.
func (f *Firmware) Update() {
	for _, u := range f.updaters {
		cv, err := u.current()
		if err != nil {
			log.Error("firmware", "unable to get current version", err)
			continue
		}
		dv := u.desired()
		log.Info("firmware", "name", u, "current", cv, "desired", dv, "update required", u.updateRequired())
		if !u.updateRequired() {
			continue
		}
		err = u.update()
		if err != nil {
			log.Error("firmware", "unable to update", err)
			continue
		}
	}
}
