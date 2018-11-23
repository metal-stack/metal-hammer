package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"

	log "github.com/inconshreveable/log15"
)

// StartSSHD will start sshd to be able to diagnose problems on the pxe bootet machine.
func StartSSHD() error {
	sshd, err := exec.LookPath("sshd")
	if err != nil {
		return fmt.Errorf("unable to locate sshd info:%v", err)
	}
	cmd := exec.Command(sshd, "-port", "22")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start sshd info:%v", err)
	}
	log.Info(fmt.Sprintf("sshd started, connect via ssh -i metal.key root@%s", getInternalIP()))
	return nil
}

func getInternalIP() string {
	var ip net.IP
	interfaces := []string{"eth0", "eth1", "eth2", "eth3", "eth4", "eth5", "eth6", "eth7", "eth8", "eth9"}
	for _, eth := range interfaces {
		itf, _ := net.InterfaceByName(eth)
		item, _ := itf.Addrs()
		for _, addr := range item {
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						ip = v.IP
					}
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	}
	return ""
}
