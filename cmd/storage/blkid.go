package storage

import (
	"os/exec"
	"strings"

	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/pkg/errors"
)

const blkidCommand = command.BlkID

// FetchBlockIDProperties use blkid to return more properties of the given partition device
func FetchBlockIDProperties(partitionDevice string) (map[string]string, error) {
	path, err := exec.LookPath(blkidCommand)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to locate program:%s in path", blkidCommand)
	}
	out, err := exec.Command(path, "-o", "export", partitionDevice).CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute %s", blkidCommand)
	}

	// output of
	// blkid /dev/sda1 -o export:
	//
	// DEVNAME=/dev/sda1
	// UUID=E562-31F0
	// TYPE=vfat
	// PARTLABEL=EFI\ System\ Partition
	// PARTUUID=5995932d-c5ba-43db-bd4b-53564510720
	//
	// we just put every key=value entry into a map

	props := make(map[string]string)
	for _, line := range strings.Split(string(out), "\n") {
		keyValue := strings.Split(line, "=")
		if len(keyValue) != 2 {
			continue
		}
		props[keyValue[0]] = keyValue[1]
	}
	return props, nil
}

// FetchBlockIDProperties use blkid to determine more properties of the partition
func (p *Partition) fetchBlockIDProperties() error {
	props, err := FetchBlockIDProperties(p.Device)
	if err != nil {
		return err
	}

	p.Properties = props
	return nil
}
