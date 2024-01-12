package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	gos "os"

	"github.com/metal-stack/metal-hammer/pkg/os"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"golang.org/x/sync/errgroup"

	"github.com/jaypipes/ghw"
)

var (
	// DiskPrefixToIgnore disks with this prefix will not be reported and wiped.
	DiskPrefixToIgnore = "ram"
)

type Disks struct {
	log *slog.Logger
}

func NewDisks(log *slog.Logger) *Disks {
	return &Disks{log: log}
}

// WipeDisks will erase all content and partitions of all existing Disks
func (d *Disks) Wipe() error {
	d.log.Info("wipe")
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks %w", err)
	}
	disks := block.Disks

	d.log.Info("wipe existing disks", "disks", disks)

	g, _ := errgroup.WithContext(context.Background())
	for _, disk := range disks {
		disk := disk
		if strings.HasPrefix(disk.Name, DiskPrefixToIgnore) {
			d.log.Info("skip because in ignorelist", "disk", disk.Name)
			continue
		}
		g.Go(func() error {
			return d.wipe(disk)
		})
	}

	err = g.Wait()
	if err != nil {
		d.log.Error("failed to wipe disk", "error", err)
	}

	return nil
}

// bs is the blocksize in bytes to be used by dd
const bs = uint64(10240)

// WipeDisk will erase all content and partitions of given existing disk.
func (d *Disks) wipe(disk *ghw.Disk) error {
	device := fmt.Sprintf("/dev/%s", disk.Name)
	bytes := disk.SizeBytes
	rotational := d.isRotational(disk.Name)
	if rotational {
		return d.insecureErase(device, bytes)
	}
	if isNVMeDisk(device) {
		return d.secureEraseNVMe(device)
	}
	return d.insecureErase(device, bytes)
}

// insecureErase will first try to format the device with discard, if this fails
// overwrite it with dd
func (d *Disks) insecureErase(device string, bytes uint64) error {
	err := d.discard(device)
	if err != nil {
		return d.wipeSlow(device, bytes)
	}
	return nil
}

func (d *Disks) discard(device string) error {
	d.log.Info("wipe", "disk", device, "message", "discard existing data")
	err := os.ExecuteCommand(command.MKFSExt4, "-F", "-E", "discard", device)
	if err != nil {
		d.log.Error("wipe", "disk", device, "message", "discard of existing data failed", "error", err)
		return err
	}

	// additionally wipe magic bytes in the first 1MiB
	err = os.ExecuteCommand(command.DD, "status=progress", "if=/dev/zero", "of="+device, "bs=1M", "count=1")
	if err != nil {
		d.log.Error("wipe", "disk", device, "message", "overwrite of the first bytes of data with dd failed", "error", err)
		return err
	}

	d.log.Info("wipe", "disk", device, "message", "finish discard of existing data")
	return nil
}

func (d *Disks) wipeSlow(device string, bytes uint64) error {
	d.log.Info("wipe", "disk", device, "message", "slow deleting of existing data")
	count := bytes / bs
	bsArg := fmt.Sprintf("bs=%d", bs)
	countArg := fmt.Sprintf("count=%d", count)
	err := os.ExecuteCommand(command.DD, "status=progress", "if=/dev/zero", "of="+device, bsArg, countArg)
	if err != nil {
		d.log.Error("wipe", "disk", device, "message", "overwrite of existing data with dd failed", "error", err)
		return err
	}
	d.log.Info("wipe", "disk", device, "message", "finish deleting of existing data")
	return nil
}

func isNVMeDisk(device string) bool {
	return strings.HasPrefix(device, "/dev/nvm")
}

// Secure erase is done via:
// nvme-cli --format --force --ses=1 /dev/nvme0n1
// see: https://github.com/linux-nvme/nvme-cli/blob/master/Documentation/nvme-format.txt
//
// TODO: configure qemu to map a disk with the nvme format:
// https://github.com/nvmecompliance/manage/blob/master/runQemu.sh
// https://github.com/arunar/nvmeqemu
func (d *Disks) secureEraseNVMe(device string) error {
	d.log.Info("wipe", "disk", device, "message", "start very fast deleting of existing data")
	err := os.ExecuteCommand(command.NVME, "--format", "--force", "--ses=1", device)
	if err != nil {
		return fmt.Errorf("unable to secure erase nvme disk %s %w", device, err)
	}
	return nil
}

func (d *Disks) isRotational(deviceName string) bool {
	sysfsRotational := fmt.Sprintf("/sys/block/%s/queue/rotational", deviceName)
	rotational, err := gos.ReadFile(sysfsRotational)
	result := true
	if err != nil {
		// defensive guess, fall back to hdd if unknown
		d.log.Warn("unable to detect if disk is rotational", "disk", deviceName, "error", err)
		return true
	}
	if strings.Contains(string(rotational), "0") {
		result = false
	}
	d.log.Debug("wipe", "disk", deviceName, "rotational", result)
	return result
}
