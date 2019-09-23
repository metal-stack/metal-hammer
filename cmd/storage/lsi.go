package storage

import (
	"encoding/json"
	"os/exec"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

type (
	// Storcli is the command to query lsi raid controllers
	Storcli struct {
		command string
	}

	// LSIController is the response of a storcli show command
	LSIController struct {
		Controllers []struct {
			CommandStatus map[string]interface{} `json:"Command Status"`
			ResponseData  map[string]interface{} `json:"Response Data"`
		} `json:"Controllers"`
	}
)

const StorCLICommand = "storcli"

// NewStorcli create a new Storcli with the default command
func NewStorcli() *Storcli {
	return &Storcli{command: StorCLICommand}
}

// Run execute ethtool
func (s *Storcli) run(args ...string) (string, error) {
	path, err := exec.LookPath(s.command)
	if err != nil {
		return "", errors.Wrapf(err, "unable to locate program:%s in path", s.command)
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()

	log.Debug("run", "command", s.command, "args", args, "output", string(output), "error", err)
	return string(output), err
}

// EnableJBOD configure all attached disks as JBOD
func (s *Storcli) EnableJBOD() error {
	if !s.controllerPresent() {
		return nil
	}

	// enable jbod on all controllers
	_, err := s.run("/call", "set", "jbod=on")
	if err != nil {
		return errors.Wrap(err, "unable set jbod on")
	}
	return nil
}

func (s *Storcli) controllerPresent() bool {

	output, err := s.run("show", "J")
	if err != nil {
		log.Error("lsi", "unable to query raidcontrollers", err)
		return false
	}

	var result LSIController
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		log.Error("lsi", "unable to unmarshal command output", err)
	}
	ctrl := result.Controllers[0]
	controllerCount := ctrl.ResponseData["Number of Controllers"]

	if controllerCount.(float64) > 0 {
		return true
	}

	return false
}
