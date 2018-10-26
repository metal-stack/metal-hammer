package cmd

import (
	"fmt"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	log "github.com/inconshreveable/log15"
)

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("metal-hammer run", "firmware", pkg.Firmware())

	err := WipeDisks(spec)
	if err != nil {
		return fmt.Errorf("wipe error: %v", err)
	}

	err = createSyslog()
	if err != nil {
		return fmt.Errorf("unable to write kernel boot message to /var/log/syslog, info:%v", err)
	}

	uuid, err := RegisterDevice(spec)
	if !spec.DevMode && err != nil {
		return fmt.Errorf("register error: %v", err)
	}

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	var device *Device
	if spec.DevMode {
		ip := "10.0.1.2/24"
		if !spec.BGPEnabled {
			ip = "dhcp"
		}
		device = &Device{
			Image: &Image{
				Url: spec.ImageURL,
			},
			Hostname:  "devmode",
			SSHPubKey: "not a valid ssh public key, can be specified during device create.",
			IP:        ip,
		}
	} else {
		device, err = Wait(spec.InstallURL, uuid)
		if err != nil {
			return fmt.Errorf("wait for installation error: %v", err)
		}
	}

	installationStart := time.Now()
	info, err := Install(device)
	if err != nil {
		return fmt.Errorf("install error: %v", err)
	}

	err = ReportInstallation(spec.ReportURL, uuid, err)
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
