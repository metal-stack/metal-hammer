package storage

import (
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/pkg/os"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
)

func ActivateRaid() error {
	log.Info("activate dmraid devices if any")
	err := os.ExecuteCommand(command.DMRaid, "-a", "y")
	if err != nil {
		log.Error("wipe", "unable to activate dmraid devices", err)
		return err
	}
	return nil
}
