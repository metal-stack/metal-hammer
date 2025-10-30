package api

import "github.com/metal-stack/metal-go/api/models"

// Bootinfo is written by the installer in the target os to tell us
// which kernel, initrd and cmdline must be used for kexec
type Bootinfo struct {
	Initrd       string `yaml:"initrd"`
	Cmdline      string `yaml:"cmdline"`
	Kernel       string `yaml:"kernel"`
	BootloaderID string `yaml:"bootloader_id"`
}

// InstallerConfig contains configuration items which are
// consumed by the install.sh / install-go of the individual target OS.
type InstallerConfig struct {
	// Hostname of the machine
	Hostname string `yaml:"hostname"`
	// Networks all networks connected to this machine
	Networks []*models.V1MachineNetwork `yaml:"networks"`
	// MachineUUID is the unique UUID for this machine, usually the board serial.
	MachineUUID string `yaml:"machineuuid"`
	// SSHPublicKey of the user
	SSHPublicKey string `yaml:"sshpublickey"`
	// Password is the password for the metal user.
	Password string `yaml:"password"`
	// Console specifies where the kernel should connect its console to.
	Console string `yaml:"console"`
	// Timestamp is the the timestamp of installer config creation.
	Timestamp string `yaml:"timestamp"`
	// Nics are the network interfaces of this machine including their neighbors.
	Nics []*models.V1MachineNic `yaml:"nics"`
	// VPN is the config for connecting machine to VPN
	VPN *models.V1MachineVPN `yaml:"vpn"`
	// Role is either firewall or machine
	Role string `yaml:"role"`
	// RaidEnabled is set to true if any raid devices are specified
	RaidEnabled bool `yaml:"raidenabled"`
	// RootUUID is the fs uuid if the root fs
	RootUUID string `yaml:"root_uuid"`
	// FirewallRules if not empty firewall rules to enforce
	FirewallRules *models.V1FirewallRules `yaml:"firewall_rules"`
	// DNSServers for the machine
	DNSServers []*models.V1DNSServer `yaml:"dns_servers"`
	// NTPServers for the machine
	NTPServers []*models.V1NTPServer `yaml:"ntp_servers"`
}

// FIXME legacy structs remove once old images are gone

type (
	// Disk is a physical Disk
	Disk struct {
		// Device the name of the disk device visible from kernel side, e.g. sda
		Device string
		// Partitions to create on this disk, order is preserved
		Partitions []Partition
	}
	Partition struct {
		Label      string
		Filesystem string
		Properties map[string]string
	}
)
