package cmd

import (
	"fmt"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	log "github.com/inconshreveable/log15"
)

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("metal-hammer run")
	log.Info("metal-hammer bootet with", "firmware", pkg.Firmware())

	err := WipeDisks(spec)
	if err != nil {
		return fmt.Errorf("wipe error: %v", err)
	}

	uuid, err := RegisterDevice(spec)
	if !spec.DevMode && err != nil {
		return fmt.Errorf("register error: %v", err)
	}

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	var device *Device
	if spec.DevMode {
		device = &Device{
			Image: &Image{
				Url: spec.ImageURL,
			},
			Hostname:  "dummy",
			SSHPubKey: "a not working key",
		}
	} else {
		device, err = Wait(spec.InstallURL, uuid)
		if err != nil {
			return fmt.Errorf("wait for installation error: %v", err)
		}
	}

	info, err := Install(device)
	if err != nil {
		return fmt.Errorf("install error: %v", err)
	}

	err = ReportInstallation(spec.ReportURL, uuid, err)
	if err != nil {
		log.Error("report install, reboot in 10sec", "error", err)
		time.Sleep(10 * time.Second)
		if !spec.DevMode {
			err = pkg.Reboot()
			if err != nil {
				log.Error("reboot", "error", err)
			}
		}
	}

	pkg.RunKexec(info)
	return nil
}
