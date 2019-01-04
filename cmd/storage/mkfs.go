package storage

import (
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/os"
	log "github.com/inconshreveable/log15"
	"strings"
)

var (
	ext4MkFsCommand  = "mkfs.ext4"
	ext3MkFsCommand  = "mkfs.ext3"
	fat32MkFsCommand = "mkfs.vfat"
	mkswapCommand    = "mkswap"
)

// MkFS create a filesystem on the Partition
func (p *Partition) MkFS() error {
	log.Info("create filesystem", "device", p.Device, "filesystem", p.Filesystem)

	mkfs, args, err := assembleMKFSCommand(p)
	if err != nil {
		return fmt.Errorf("mkfs failed with error:%v", err)
	}
	err = os.ExecuteCommand(mkfs, args...)
	if err != nil {
		return fmt.Errorf("mkfs failed with error:%v", err)
	}

	return nil
}

func assembleMKFSCommand(p *Partition) (string, []string, error) {
	mkfs := ""
	var args []string
	switch p.Filesystem {
	case EXT4:
		mkfs = ext4MkFsCommand
		args = append(args, "-v", "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case EXT3:
		mkfs = ext3MkFsCommand
		args = append(args, "-v", "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case FAT32, VFAT:
		mkfs = fat32MkFsCommand
		args = append(args, "-v", "-F", "32")
		if p.Label != "" {
			args = append(args, "-n", strings.ToUpper(p.Label))
		}
	case SWAP:
		mkfs = mkswapCommand
		args = append(args, "-f")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	default:
		return "", nil, fmt.Errorf("unsupported filesystem format: %q", p.Filesystem)
	}
	args = append(args, p.Device)
	return mkfs, args, nil
}
