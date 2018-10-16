package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	log "github.com/inconshreveable/log15"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// downloadFile will download from a source url to a local file dest.
// It's efficient because it will write as it downloads
// and not load the whole file into memory.
func download(source, dest string) error {
	log.Info("download", "from", source, "to", dest)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("download of %s did not work, statuscode was: %d", source, resp.StatusCode)
	}

	fileSize := resp.ContentLength

	bar := pb.New64(fileSize).SetUnits(pb.U_BYTES)
	bar.SetWidth(80)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	reader := bar.NewProxyReader(resp.Body)
	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	return nil
}

// small helper to execute a command, redirect stdout/stderr.
func executeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
