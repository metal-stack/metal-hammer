package command

import (
	"fmt"
	"os/exec"
)

const (
	BlkID    = "blkid"
	DD       = "dd"
	MDADM    = "mdadm"
	LVM      = "lvm"
	Ethtool  = "ethtool"
	HDParm   = "hdparm"
	IPMITool = "ipmitool"
	MKFSExt3 = "mkfs.ext3"
	MKFSExt4 = "mkfs.ext4"
	MKFSVFat = "mkfs.vfat"
	MKSwap   = "mkswap"
	NVME     = "nvme"
	SGDisk   = "sgdisk"
	SSHD     = "sshd"
	SUM      = "sum"
	WIPEFS   = "wipefs"
)

var commands = []string{
	BlkID,
	DD,
	MDADM,
	LVM,
	Ethtool,
	HDParm,
	IPMITool,
	MKFSExt3,
	MKFSExt4,
	MKFSVFat,
	MKSwap,
	NVME,
	SGDisk,
	SSHD,
	SUM,
	WIPEFS,
}

// CommandsExist check that all required binaries are installed in the initrd.
func CommandsExist() error {
	missingCommands := []string{}
	for _, command := range commands {
		_, err := exec.LookPath(command)
		if err != nil {
			missingCommands = append(missingCommands, command)
		}
	}
	if len(missingCommands) > 0 {
		return fmt.Errorf("unable to locate:%s in path", missingCommands)
	}
	return nil
}
