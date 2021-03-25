package storage

import (
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/pkg/os"
)

func ActivateRaid() error {
	log.Info("activate sw raid devices if any")
	// err := os.ExecuteCommand(command.MDADM, "-A", "-s")
	err := os.ExecuteCommand("dmraid", "-a", "y")
	if err != nil {
		log.Error("wipe", "unable to activate sw raid devices", err)
		return err
	}
	return nil
}
