package cmd

import (
	"fmt"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/models"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/inconshreveable/log15"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Client *device.Client
	Spec   *Specification
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("metal-hammer run", "firmware", pkg.Firmware())

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := device.New(transport, strfmt.Default)

	hammer := &Hammer{
		Client: client,
		Spec:   spec,
	}

	err := hammer.WipeDisks()
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
	var device *models.ModelsMetalDevice
	if spec.DevMode {
		cidr := "10.0.1.2/24"
		if !spec.BGPEnabled {
			cidr = "dhcp"
		}
		hostname := "devmode"
		sshkey := "not a valid ssh public key, can be specified during device create."
		device = &models.ModelsMetalDevice{
			Image: &models.ModelsMetalImage{
				URL: &spec.ImageURL,
			},
			Hostname:  &hostname,
			SSHPubKey: &sshkey,
			Cidr:      &cidr,
		}
	} else {
		device, err = hammer.Wait(uuid)
		if err != nil {
			return fmt.Errorf("wait for installation error: %v", err)
		}
	}

	installationStart := time.Now()
	info, err := Install(device)
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
