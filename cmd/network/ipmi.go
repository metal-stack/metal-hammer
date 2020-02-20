package network

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/ipmi"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const defaultIpmiPort = "623"

// defaultIpmiUser the name of the user created by metal in the ipmi config
const defaultIpmiUser = "metal"

// defaultIpmiUserID the id of the user created by metal in the ipmi config
const defaultIpmiUserID = "10"

// IPMI configuration and
func readIPMIDetails(eth0Mac string) (*models.ModelsV1MachineIPMI, error) {
	config := ipmi.LanConfig{}
	var fru models.ModelsV1MachineFru
	i := ipmi.New()
	var pw string
	var user string
	var bmcInfo ipmi.BMCInfo
	if i.DevicePresent() {
		log.Info("ipmi details from bmc")
		user = defaultIpmiUser
		// FIXME userid should be verified if available
		var err error
		pw, err = i.CreateUser(user, defaultIpmiUserID, ipmi.Administrator)
		if err != nil {
			return nil, errors.Wrap(err, "ipmi create user failed")
		}
		config, err = i.GetLanConfig()
		if err != nil {
			return nil, errors.Wrap(err, "unable to read ipmi lan configuration")
		}
		log.Debug("register", "ipmi lanconfig", config)
		config.IP = config.IP + ":" + defaultIpmiPort
		f, err := i.GetFru()
		if err != nil {
			return nil, errors.Wrap(err, "unable to read ipmi fru configuration")
		}
		bmcInfo, err = i.GetBMCInfo()
		if err != nil {
			return nil, errors.Wrap(err, "unable to read ipmi bmc info configuration")
		}
		fru = models.ModelsV1MachineFru{
			ChassisPartNumber:   f.ChassisPartNumber,
			ChassisPartSerial:   f.ChassisPartSerial,
			BoardMfg:            f.BoardMfg,
			BoardMfgSerial:      f.BoardMfgSerial,
			BoardPartNumber:     f.BoardPartNumber,
			ProductManufacturer: f.ProductManufacturer,
			ProductPartNumber:   f.ProductPartNumber,
			ProductSerial:       f.ProductSerial,
		}

	} else {
		log.Info("ipmi details faked")

		if len(eth0Mac) == 0 {
			eth0Mac = "00:00:00:00:00:00"
		}

		macParts := strings.Split(eth0Mac, ":")
		lastOctet := macParts[len(macParts)-1]
		port, err := strconv.ParseUint(lastOctet, 16, 32)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse last octet of eth0 mac to a integer")
		}

		const baseIPMIPort = 6230
		// Fixed IP of vagrant environment gateway
		config.IP = fmt.Sprintf("192.168.121.1:%d", baseIPMIPort+port)
		config.Mac = "00:00:00:00:00:00"
		pw = "vagrant"
		user = "vagrant"
	}

	intf := "lanplus"
	details := &models.ModelsV1MachineIPMI{
		Address:    &config.IP,
		Mac:        &config.Mac,
		Password:   &pw,
		User:       &user,
		Interface:  &intf,
		Fru:        &fru,
		Bmcversion: &bmcInfo.FirmwareRevision,
	}

	return details, nil
}
