package cmd

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/inconshreveable/log15"
)

var (
	imgCommand = "/bin/img"
)

// Install a given image to the disk by using genuinetools/img
func Install(image string) error {
	err := pull(image)
	if err != nil {
		return err
	}
	err = burn(image)
	if err != nil {
		return err
	}
	return nil
}

// pull a image by calling genuinetools/img pull
func pull(image string) error {
	cmd := exec.Command(imgCommand, "pull", image)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to pull image %s error message: %v error: %v", image, string(output), err)
	}
	log.Debug("pull image", "output", output, "image", image)
	return nil
}

// burn a image by calling genuinetools/img unpack to a specific directory
func burn(image string) error {
	cmd := exec.Command(imgCommand, "unpack", image)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to burn image %s error message: %v error: %v", image, string(output), err)
	}
	log.Debug("burn image", "output", output, "image", image)
	return nil
}
