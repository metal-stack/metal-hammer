package storage

import (
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/os"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

var (
	sgdiskCommand = "sgdisk"
)

// Partition a Disk
func (disk Disk) Partition() error {
	log.Info("partition disk", "disk", disk)

	err := os.ExecuteCommand(sgdiskCommand, "-Z", disk.Device)
	if err != nil {
		log.Error("sgdisk zapping existing partitions failed, ignoring...", "error", err)
	}

	args := assembleSGDiskCommand(disk)
	log.Info("sgdisk create partitions", "command", args)
	err = os.ExecuteCommand(sgdiskCommand, args...)
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err)
		return errors.Wrapf(err, "unable to create partitions on %s", disk)
	}

	return nil
}

func assembleSGDiskCommand(disk Disk) []string {
	args := make([]string, 0)
	for _, p := range disk.Partitions {
		size := fmt.Sprintf("%dM", p.Size)
		if p.Size == -1 {
			size = "0"
		}
		args = append(args, fmt.Sprintf("-n=%d:0:%s", p.Number, size))
		args = append(args, fmt.Sprintf("-c=%d:%s", p.Number, p.Label))
		args = append(args, fmt.Sprintf("-t=%d:%s", p.Number, p.GPTType))
		if p.GPTGuid != "" {
			args = append(args, fmt.Sprintf("-u=%d:%s", p.Number, p.GPTGuid))
		}
	}

	args = append(args, disk.Device)
	return args
}
