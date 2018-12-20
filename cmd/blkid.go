package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

const blkidCommand = "blkid"

func (p *Partition) fetchBlockIDProperties() error {

	path, err := exec.LookPath(blkidCommand)
	if err != nil {
		return fmt.Errorf("unable to locate program:%s in path info:%v", blkidCommand, err)
	}
	out, err := exec.Command(path, "-o", "export", p.Device).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to execute %s error:%v", blkidCommand, err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		keyValue := strings.Split(line, "=")
		if len(keyValue) != 2 {
			continue
		}
		p.Properties[keyValue[0]] = keyValue[1]
	}
	return nil
}
