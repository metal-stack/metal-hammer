package firmware

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
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
	log      *zap.SugaredLogger
}

// New create a new Firmware manager with all Updaters.
func New(log *zap.SugaredLogger) *Firmware {

	_ = raidcontroller{
		name:           "lsi3108",
		desiredVersion: "4.680.00-8290",
		log:            log,
	}
	_ = intel{
		name:           "intel nics",
		desiredVersion: "6.8",
		log:            log,
	}
	return &Firmware{
		updaters: []updater{},
		log:      log,
	}
}

// Update run updates for all firmwares found.
func (f *Firmware) Update() {
	for _, u := range f.updaters {
		cv, err := u.current()
		if err != nil {
			f.log.Errorw("firmware", "unable to get current version", err)
			continue
		}
		dv := u.desired()
		f.log.Infow("firmware", "name", u, "current", cv, "desired", dv, "update required", u.updateRequired())
		if !u.updateRequired() {
			continue
		}
		err = u.update()
		if err != nil {
			f.log.Errorw("firmware", "unable to update", err)
			continue
		}
	}
}

// Run execute a comand with arguments, returns output and error
func run(log *zap.SugaredLogger, command string, args ...string) (string, error) {
	path, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("unable to locate program:%s in path %w", command, err)
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()

	log.Debugw("run", "command", command, "args", args, "output", string(output), "error", err)
	return string(output), err
}
