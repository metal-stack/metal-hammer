package storage

import (
	"fmt"
	"os/exec"
	"strings"
)

const blkidCommand = "blkid"

// FetchBlockIDProperties use blkid to determine more properties of the partition
func (p *Partition) fetchBlockIDProperties() error {

	path, err := exec.LookPath(blkidCommand)
	if err != nil {
		return fmt.Errorf("unable to locate program:%s in path info:%v", blkidCommand, err)
	}
	out, err := exec.Command(path, "-o", "export", p.Device).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to execute %s error:%v", blkidCommand, err)
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

	for _, line := range strings.Split(string(out), "\n") {
		keyValue := strings.Split(line, "=")
		if len(keyValue) != 2 {
			continue
		}
		p.Properties[keyValue[0]] = keyValue[1]
	}
	return nil
}
