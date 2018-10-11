package cmd

import (
	"fmt"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

// WipeDisks will erase all content and partitions of all existing Disks
func WipeDisks() error {
	log.Info("wipe all disks")
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks: %v", err)
	}
	for _, disk := range block.Disks {
		log.Info("TODO wipe disk", "disk", disk)

		diskDevice := fmt.Sprintf("/dev/%s", disk.Name)
		log.Info("sgdisk zap all existing partitions", "disk", diskDevice)
		err := executeCommand(sgdiskCommand, "-Z", diskDevice)
		if err != nil {
			log.Error("sgdisk zap all existing partitions failed", "error", err)
		}
	}
	return nil
}
