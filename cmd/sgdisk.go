package cmd

import (
	"fmt"

	log "github.com/inconshreveable/log15"
)

var (
	sgdiskCommand = "sgdisk"
)

func partition(disk Disk) error {
	log.Info("partition disk", "disk", disk)

	err := executeCommand(sgdiskCommand, "-Z", disk.Device)
	if err != nil {
		log.Error("sgdisk zapping existing partitions failed, ignoring...", "error", err)
	}

	args := assembleSGDiskCommand(disk)
	log.Info("sgdisk create partitions", "command", args)
	err = executeCommand(sgdiskCommand, args...)
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err)
		return fmt.Errorf("unable to create partitions on %s error:%v", disk, err)
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
