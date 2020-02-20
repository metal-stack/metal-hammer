package cmd

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

// fetchMachine requests the machine data of given machine ID
func (h *Hammer) fetchMachine(machineID string) (*models.ModelsV1MachineResponse, error) {
	params := machine.NewFindMachineParams()
	params.SetID(machineID)
	resp, err := h.Client.FindMachine(params)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (h *Hammer) reinstall(m *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended, eventEmitter *event.EventEmitter) (*event.EventEmitter, error) {
	// wipe only the disk that has the OS installed on one of its partitions, keep all other disks untouched
	primaryDiskName := *m.Allocation.Reinstallation.PrimaryDisk
	var primaryBlockDevice *models.ModelsV1MachineBlockDevice
	for _, d := range m.Hardware.Disks {
		if *d.Name == primaryDiskName {
			primaryBlockDevice = d
			break
		}
	}
	if primaryBlockDevice == nil {
		return eventEmitter, fmt.Errorf("unable to find primary disk %s", primaryDiskName)
	}
	if h.Disk.Device != primaryDiskName {
		return eventEmitter, fmt.Errorf("new image OS %s is not compatible to existing OS", m.Allocation.Image.Name)
	}
	osPartitionName := *m.Allocation.Reinstallation.OsPartition
	var osPartition *models.ModelsV1MachineDiskPartition
	for _, p := range primaryBlockDevice.Partitions {
		if *p.Device == osPartitionName {
			osPartition = p
			break
		}
	}
	if osPartition == nil {
		return eventEmitter, fmt.Errorf("unable to find partition %s on primary disk %s that has OS installed", osPartitionName, primaryDiskName)
	}

	block, err := ghw.Block()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "unable to gather disks")
	}
	var primaryDisk *ghw.Disk
	for _, d := range block.Disks {
		if d.Name == primaryDiskName {
			primaryDisk = d
			break
		}
	}
	if primaryDisk == nil {
		return eventEmitter, errors.Wrapf(err, "unable to find primary disk %s", primaryDiskName)
	}
	var primaryPartition *storage.Partition
	for _, p := range h.Disk.Partitions {
		if p.Device == osPartitionName {
			primaryPartition = p
			break
		}
	}
	if primaryPartition == nil {
		return eventEmitter, errors.Wrapf(err, "unable to find primary partition %s", osPartitionName)
	}

	err = h.verifyOS(osPartition)
	if err != nil {
		log.Error("reinstallation cancelled => boot into existing OS...", "primary disk", primaryDiskName, "partition device", osPartition.Device, "error", err)
		h.EventEmitter.Emit(event.ProvisioningEventReinstallAborted, fmt.Sprintf("Reinstallation aborted: %s", err.Error()))
		err = h.bootHD(h.Spec.MachineUUID)
		return eventEmitter, err
	}

	err = storage.WipeDisk(primaryDisk)
	if err != nil {
		return eventEmitter, errors.Wrap(err, "wipe")
	}

	return h.installImage(eventEmitter, m, nics, false)
}

// verifyOS checks if there is indeed an OS installed on given partition.
// For this it mounts that partition to /rootfs/verify-os and checks if /rootfs/verify-os/etc/metal is present.
func (h *Hammer) verifyOS(part *models.ModelsV1MachineDiskPartition) error {
	err := h.Disk.Partition()
	if err != nil {
		return err
	}

	p := h.Disk.Partitions[0]
	err = p.MkFS()
	if err != nil {
		return err
	}

	props, err := storage.FetchBlockIDProperties(*part.Device)
	if err != nil {
		log.Error("unable to fetch blockid properties of partition device", "partition device", *part.Device, "error", err)
		return errors.Wrapf(err, "unable to fetch blockid properties of partition device %s", *part.Device)
	}

	fstype := props["TYPE"]
	data := ""
	flags := uintptr(syscall.MS_BIND)

	mountPoint := filepath.Join(h.ChrootPrefix, "verify-os")
	log.Info("mount partition device", "source", *part.Device, "target", mountPoint, "fstype", fstype, "flags", flags, "data", data)
	err = syscall.Mount(*part.Device, mountPoint, fstype, flags, data)
	if err != nil {
		log.Error("unable to mount partition device", "path", mountPoint, "error", err)
		return errors.Wrapf(err, "mounting partition device %s to %s failed", *part.Device, mountPoint)
	}

	_, etcMetal := os.Stat(fmt.Sprintf("%s/etc/metal", mountPoint))
	log.Info("unmount partition device", "mountpoint", mountPoint)
	err = syscall.Unmount(mountPoint, syscall.MNT_FORCE)
	if err != nil {
		log.Error("unable to unmount partition device... ignored", "path", mountPoint, "error", err)
	}
	if os.IsNotExist(etcMetal) {
		log.Error("no OS found on partition device", "partition device", *part.Device)
		return fmt.Errorf("no OS found on partition device %s", *part.Device)
	}

	return nil
}

// bootHD instructs metal-core to change boot order of this machine to HD.
// It then kexecs into existing kernel.
func (h *Hammer) bootHD(machineID string) error {
	cbo := &models.DomainChangeBootOrder{
		Hd: true,
	}
	params := machine.NewAbortReinstallParams()
	params.SetID(machineID)
	params.SetBody(cbo)
	_, err := h.Client.AbortReinstall(params)
	if err != nil {
		return err
	}

	info, err := kernel.ReadBootinfo(path.Join(h.ChrootPrefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return err
	}

	return kernel.RunKexec(info)
}
