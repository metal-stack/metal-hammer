package cmd

import (
	"fmt"
	"sync"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

// TODO, check with hdparm -i for Self Encrypting Device, sample output will look like:
// Security:
//         Master password revision code = 65534
//                 supported
//         not     enabled
//         not     locked
//                 frozen
//         not     expired: security count
//                 supported: enhanced erase
//         6min for SECURITY ERASE UNIT. 32min for ENHANCED SECURITY ERASE UNIT.

// WipeDisks will erase all content and partitions of all existing Disks
func WipeDisks(spec *Specification) error {
	log.Info("wipe all disks", "devmode", spec.DevMode)
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks: %v", err)
	}
	disks := block.Disks

	log.Info("wipe existing disks", "disks", disks)

	wipeErrors := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(disks))
	for _, disk := range disks {
		go func(disk *ghw.Disk) {
			defer wg.Done()
			device := fmt.Sprintf("/dev/%s", disk.Name)
			bytes := disk.SizeBytes

			err := wipeDisk(device, bytes)
			if err != nil {
				wipeErrors <- err
			}
		}(disk)
	}

	go func() {
		for e := range wipeErrors {
			log.Error("failed to wipe disk", "error", e)
		}
	}()
	wg.Wait()

	return nil
}

const bs = uint64(10240)

func wipeDisk(device string, bytes uint64) error {
	log.Info("start deleting of existing data on", "disk", device)
	count := bytes / bs
	bsArg := fmt.Sprintf("bs=%d", bs)
	countArg := fmt.Sprintf("count=%d", count)
	err := executeCommand("/bbin/dd", "if=/dev/zero", "of="+device, bsArg, countArg)
	if err != nil {
		log.Error("overwrite of existing data with dd failed", "disk", device, "error", err)
		return err
	}
	log.Info("finish deleting of existing data on", "disk", device)
	return nil
}
