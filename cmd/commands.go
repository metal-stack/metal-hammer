package cmd

import (
	"fmt"
	"os/exec"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/storage"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
)

var commands = []string{
	ipmi.Command,
	network.EthtoolCommand,
	// storage.StorCLICommand, actually disabled.
	storage.BlkidCommand,
	storage.SgdiskCommand,
	storage.Ext3MkFsCommand,
	storage.Ext4MkFsCommand,
	storage.Fat32MkFsCommand,
	storage.MkswapCommand,
	storage.DDCommand,
	storage.NvmeCommand,
	storage.HdparmCommand,
	sshdCommand,
}

// checkAllCommandsExist check that all required binaries are installed in the initrd.
func checkAllCommandsExist() error {
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
