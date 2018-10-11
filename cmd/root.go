package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/inconshreveable/log15"

	"github.com/u-root/u-root/pkg/kexec"
	"golang.org/x/sys/unix"
)

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("metal-hammer run")
	firmware := bootedWith()
	log.Info("metal-hammer bootet with", "firmware", firmware)

	err := WipeDisks()
	if err != nil {
		log.Error("register device", "error", err)
	}

	uuid, err := RegisterDevice(spec)
	if err != nil {
		log.Error("register device", "error", err)
	}

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	var device *Device
	devMode := false
	if spec.ImageURL != "" {
		device = &Device{
			Image: &Image{
				Url: spec.ImageURL,
			},
			Hostname:  "dummy",
			SSHPubKey: "a not working key",
		}
		devMode = true
	} else {
		device, err = waitForInstall(spec.InstallURL, uuid)
		if err != nil {
			log.Error("wait for install", "error", err)
		}
	}

	info, err := Install(device)
	if err != nil {
		log.Error("install", "error", err)
	}

	err = ReportInstallation(spec.ReportURL, uuid, err)
	if err != nil {
		log.Error("report install, reboot in 10sec", "error", err)
		time.Sleep(10 * time.Second)
		if !devMode {
			reboot()
		}
	}

	runKexec(info)
	return nil
}

func bootedWith() string {
	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		return "bios"
	}
	return "efi"
}

func waitForInstall(url, uuid string) (*Device, error) {
	log.Info("waiting for install, long polling", "uuid", uuid)
	e := fmt.Sprintf("%v/%v", url, uuid)

	var resp *http.Response
	var err error
	for {
		resp, err = http.Get(e)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Debug("waiting for install failed", "error", err)
		} else {
			break
		}
		log.Debug("Retrying...")
		time.Sleep(2 * time.Second)
	}

	deviceJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed with: %v", err)
	}

	var device Device
	err = json.Unmarshal(deviceJSON, &device)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response with error: %v", err)
	}
	log.Debug("stopped waiting and got", "device", device)

	return &device, nil
}

func reboot() {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
		log.Error("unable to reboot", "error", err.Error())
	}
}

func runKexec(info *bootinfo) {
	kernel, err := os.OpenFile(info.Kernel, os.O_RDONLY, 0)
	if err != nil {
		log.Error("could not open", "kernel", info.Kernel, "error", err)
		return
	}
	defer kernel.Close()

	ramfs, err := os.OpenFile(info.Initrd, os.O_RDONLY, 0)
	if err != nil {
		log.Error("could not open", "initrd", info.Initrd, "error", err)
		return
	}
	defer ramfs.Close()

	if err := kexec.FileLoad(kernel, ramfs, info.Cmdline); err != nil {
		log.Error("could not execute kexec load", "info", info, "error", err)
	}

	err = kexec.Reboot()
	if err != nil {
		log.Error("could not fire kexec reboot", "info", info, "error", err)
	}
}
