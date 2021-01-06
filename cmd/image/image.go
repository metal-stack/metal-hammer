package image

import (
	"fmt"
	"path/filepath"

	log "github.com/inconshreveable/log15"
	"github.com/klauspost/compress/zstd"
	"github.com/mholt/archiver"
	lz4 "github.com/pierrec/lz4"

	pb "gopkg.in/cheggaaa/pb.v1"

	"crypto/md5"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Pull a image from s3
func Pull(image, destination string) error {
	log.Info("pull image", "image", image)
	md5destination := destination + ".md5"
	md5file := image + ".md5"
	err := download(image, destination)
	if err != nil {
		return errors.Wrapf(err, "unable to pull image %s", image)
	}
	err = download(md5file, md5destination)
	defer os.Remove(md5destination)
	if err != nil {
		return errors.Wrapf(err, "unable to pull md5 %s", md5file)
	}
	log.Info("check md5")
	matches, err := checkMD5(destination, md5destination)
	if err != nil || !matches {
		return errors.New("md5sum mismatch")
	}

	log.Info("pull image done", "image", image)
	return nil
}

// Burn a image pulling a tarball and unpack to a specific directory
func Burn(prefix, image, source string) error {
	log.Info("burn image", "image", image)
	begin := time.Now()

	file, err := os.Open(source)
	if err != nil {
		return errors.Wrapf(err, "%s: failed to open archive", source)

	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.Wrapf(err, "unable to stat %s", source)
	}

	suffix := filepath.Ext(image)
	var decompressor io.ReadCloser
	switch suffix {
	case ".lz4":
		lz4Reader := lz4.NewReader(file)
		log.Info("lz4", "size", lz4Reader.Header.Size)
		decompressor = ioutil.NopCloser(lz4Reader)
	case ".zstd":
		zstdReader, err := zstd.NewReader(file)
		if err != nil {
			return errors.Errorf("unable to create a zstd reader of image:%s", image)
		}
		defer zstdReader.Close()
		decompressor = ioutil.NopCloser(zstdReader)
	default:
		return errors.Errorf("unsupported image compression format of image:%s", image)
	}

	// wild guess for lz4 compression ratio
	// lz4/zstd is a stream format and therefore the
	// final size cannot be calculated upfront
	csize := stat.Size() * 2
	defer decompressor.Close()

	bar := pb.New64(csize).SetUnits(pb.U_BYTES)
	bar.Start()
	bar.SetWidth(80)
	bar.ShowSpeed = true

	reader := bar.NewProxyReader(decompressor)

	err = archiver.Tar.Read(reader, prefix)
	if err != nil {
		return errors.Wrapf(err, "unable to burn image %s", source)
	}

	bar.Finish()

	err = os.Remove(source)
	if err != nil {
		log.Warn("burn image unable to remove image source", "error", err)
	}

	log.Info("burn took", "duration", time.Since(begin))
	return nil
}

// checkMD5 check the md5 signature of file with the md5sum given in the md5file.
// the content of the md5file must be in the form:
// <md5sum> filename
// this is the same format as create by the "md5sum" unix command
func checkMD5(file, md5file string) (bool, error) {
	md5fileContent, err := ioutil.ReadFile(md5file)
	if err != nil {
		return false, errors.Wrapf(err, "unable to read md5sum file %s", md5file)
	}
	expectedMD5 := strings.Split(string(md5fileContent), " ")[0]

	f, err := os.Open(file)
	if err != nil {
		return false, errors.Wrapf(err, "unable to read file: %s", file)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, errors.Wrapf(err, "unable to calculate md5sum of file: %s", file)
	}
	sourceMD5 := fmt.Sprintf("%x", h.Sum(nil))
	log.Info("check md5", "source md5", sourceMD5, "expected md5", expectedMD5)
	if sourceMD5 != expectedMD5 {
		return false, errors.Errorf("source md5:%s expected md5:%s", sourceMD5, expectedMD5)
	}
	return true, nil
}

// downloadFile will download from a source url to a local file dest.
// It's efficient because it will write as it downloads
// and not load the whole file into memory.
func download(source, dest string) error {
	log.Info("download", "from", source, "to", dest)
	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "unable to create destination %s", dest)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("download of %s did not work, statuscode was: %d", source, resp.StatusCode)
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
