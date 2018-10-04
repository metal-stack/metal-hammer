package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/inconshreveable/log15"

	"golang.org/x/sys/unix"
)

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) error {
	log.Info("discover run")

	uuid, err := RegisterDevice(spec)
	if err != nil {
		log.Error("register device", "error", err)
	}

	url, err := waitForInstall(spec.InstallURL, uuid)
	if err != nil {
		log.Error("wait for install", "error", err)
	}

	err = Install(url)
	if err != nil {
		log.Error("install", "error", err)
	}

	err = reportInstallation()
	if err != nil {
		log.Error("report install", "error", err)
	}

	reboot()
	return nil
}

func waitForInstall(url, uuid string) (string, error) {
	log.Info("waiting for install", "uuid", uuid)

	e := fmt.Sprintf("%v/%v", url, uuid)
	resp, err := http.Get(e)
	if err != nil {
		return "", fmt.Errorf("waiting for install failed with: %v", err)
	}
	defer resp.Body.Close()
	imgURL, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("reading response failed with: %v", err)
	}
	return string(imgURL), nil
}

func reportInstallation() error {
	log.Info("report image installation status back")
	return nil
}

func reboot() {
	if err := unix.Reboot(int(unix.LINUX_REBOOT_CMD_RESTART)); err != nil {
		log.Error("unable to reboot", "error", err.Error())
	}
}
