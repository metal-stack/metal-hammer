package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

var (
	hdparmCommand = "hdparm"
	nvmeCommand   = "nvme"
)

// WipeDisks will erase all content and partitions of all existing Disks
func (h *Hammer) WipeDisks() error {
	log.Info("wipe all disks", "devmode", h.Spec.DevMode)
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

			err := wipe(device, bytes)
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

// bs is the blocksize in bytes to be used by dd
const bs = uint64(10240)

func wipe(device string, bytes uint64) error {
	if isSEDAvailable(device) {
		return secureErase(device)
	} else if isNVMeDisk(device) {
		return secureEraseNVMe(device)
	}
	return wipeSlow(device, bytes)
}

func wipeSlow(device string, bytes uint64) error {
	log.Info("start slow deleting of existing data on", "disk", device)
	count := bytes / bs
	bsArg := fmt.Sprintf("bs=%d", bs)
	countArg := fmt.Sprintf("count=%d", count)
	err := executeCommand("/bbin/dd", "status=progress", "if=/dev/zero", "of="+device, bsArg, countArg)
	if err != nil {
		log.Error("overwrite of existing data with dd failed", "disk", device, "error", err)
		return err
	}
	log.Info("finish deleting of existing data on", "disk", device)
	return nil
}

// isSEDAvailable check the disk if it is a Self Encryption Device
// check with hdparm -i for Self Encrypting Device, sample output will look like:
// Security:
//         Master password revision code = 65534
//                 supported
//         not     enabled
//         not     locked
//                 frozen
//         not     expired: security count
//                 supported: enhanced erase
//         6min for SECURITY ERASE UNIT. 32min for ENHANCED SECURITY ERASE UNIT.
// explanation is here: https://wiki.ubuntuusers.de/SSD/Secure-Erase/
func isSEDAvailable(device string) bool {
	cmd := exec.Command(hdparmCommand, "-i", device)
	output, err := cmd.Output()
	if err != nil {
		log.Error("error executing hdparm", "error", err)
		return false
	}
	hdparmOutput := string(output)
	if strings.Contains(hdparmOutput, "supported: enhanced erase") {

		scanner := bufio.NewScanner(strings.NewReader(hdparmOutput))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "frozen") && !strings.Contains(line, "not") {
				return false
			}
			if strings.Contains(line, "supported: enhanced erase") && strings.Contains(line, "not") {
				return false
			}
		}
		return true
	}
	return false
}

func isNVMeDisk(device string) bool {
	if strings.HasPrefix(device, "/dev/nvm") {
		return true
	}
	return false
}

// Secure erase is done via:
// nvme-cli --format --ses=1 /dev/nvme0n1
// see: https://github.com/linux-nvme/nvme-cli/blob/master/Documentation/nvme-format.txt
func secureEraseNVMe(device string) error {
	log.Info("start very fast deleting of existing data on", "disk", device)
	err := executeCommand(nvmeCommand, "--format", "--ses=1", device)
	if err != nil {
		return fmt.Errorf("unable to secure erase nvme disk %s error:%v", device, err)
	}
	return nil
}

func secureErase(device string) error {
	log.Info("start fast deleting of existing data on", "disk", device)
	// hdparm --user-master u --security-set-pass GEHEIM /dev/sda
	// FIXME random password
	password := "GEHEIM"
	err := executeCommand(hdparmCommand, "--user-master", "u", "--security-set-pass", password, device)
	if err != nil {
		return fmt.Errorf("unable to secure erase disk %s error: %v", device, err)
	}
	return nil
}
