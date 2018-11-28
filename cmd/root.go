package cmd

import (
	"fmt"
	"time"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/inconshreveable/log15"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Client *device.Client
	Spec   *Specification
	// IPAddress is the ip of the eth0 interface during installation
	IPAddress string
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("metal-hammer run", "firmware", pkg.Firmware())

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := device.New(transport, strfmt.Default)

	hammer := &Hammer{
		Client:    client,
		Spec:      spec,
		IPAddress: getInternalIP(),
	}

	err := hammer.UpAllInterfaces()
	if err != nil {
		return fmt.Errorf("interfaces error: %v", err)
	}

	err = hammer.EnsureUEFI()
	if err != nil {
		return fmt.Errorf("uefi error: %v", err)
	}

	err = hammer.WipeDisks()
	if err != nil {
		return fmt.Errorf("wipe error: %v", err)
	}

	err = createSyslog()
	if err != nil {
		return fmt.Errorf("unable to write kernel boot message to /var/log/syslog, info:%v", err)
	}

	uuid, err := hammer.RegisterDevice()
	if !spec.DevMode && err != nil {
		return fmt.Errorf("register error: %v", err)
	}

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	var deviceWithToken *models.ModelsMetalDeviceWithPhoneHomeToken
	if spec.DevMode {
		cidr := "10.0.1.2/24"
		if spec.Cidr != "" {
			cidr = spec.Cidr
		}

		if !spec.BGPEnabled {
			cidr = "dhcp"
		}
		hostname := "devmode"
		sshkeys := []string{"not a valid ssh public key, can be specified during device create.", "second public key"}
		fakeToken := "JWT"
		deviceWithToken = &models.ModelsMetalDeviceWithPhoneHomeToken{
			Device: &models.ModelsMetalDevice{
				Allocation: &models.ModelsMetalDeviceAllocation{
					Image: &models.ModelsMetalImage{
						URL: &spec.ImageURL,
					},
					Hostname:   &hostname,
					SSHPubKeys: sshkeys,
					Cidr:       &cidr,
				},
			},
			PhoneHomeToken: &fakeToken,
		}
	} else {
		deviceWithToken, err = hammer.Wait(uuid)
		if err != nil {
			return fmt.Errorf("wait for installation error: %v", err)
		}
	}

	installationStart := time.Now()
	info, err := hammer.Install(deviceWithToken)
	if err != nil {
		return fmt.Errorf("install error: %v", err)
	}

	err = hammer.ReportInstallation(uuid, err)
	if err != nil {
		wait := 10 * time.Second
		log.Error("report installation failed", "reboot in", wait, "error", err)
		time.Sleep(wait)
		if !spec.DevMode {
			err = pkg.Reboot()
			if err != nil {
				log.Error("reboot", "error", err)
			}
		}
	}

	log.Info("installation", "took", time.Since(installationStart))
	return pkg.RunKexec(info)
}
