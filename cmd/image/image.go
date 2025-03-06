package image

import (
	"crypto/sha512"
	"fmt"
	"hash"
	"log/slog"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/mholt/archiver"
	lz4 "github.com/pierrec/lz4/v4"

	//nolint:gosec
	"crypto/md5"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	HashMD5    = "md5"
	HashSHA512 = "sha512"
)

type Image struct {
	log *slog.Logger
}

func NewImage(log *slog.Logger) *Image {
	return &Image{log: log}
}

// Pull a image from s3
func (i *Image) Pull(image, destination string) error {
	var (
		sha512destination = destination + ".sha512sum"
		sha512file        = image + ".sha512sum"
		md5destination    = destination + ".md5"
		md5file           = image + ".md5"
	)

	i.log.Info("pull image", "image", image)
	err := i.download(image, destination)
	if err != nil {
		return fmt.Errorf("unable to pull image %s %w", image, err)
	}

	err = i.download(sha512file, sha512destination)
	defer os.Remove(sha512destination)
	if err != nil {
		i.log.Info("unable to process sha512 file, trying with md5", "error", err)
		err = i.download(md5file, md5destination)
		defer os.Remove(md5destination)
		if err != nil {
			return fmt.Errorf("unable to pull hash file %s %w", md5file, err)
		}
		matches, err := i.checkHash(destination, md5destination, HashMD5)
		if err != nil || !matches {
			return fmt.Errorf("md5 mismatch, matches: %v with error: %w", matches, err)
		}
	} else {
		i.log.Info("check sha512")
		matches, err := i.checkHash(destination, sha512destination, HashSHA512)
		if err != nil || !matches {
			return fmt.Errorf("sha512 mismatch, matches: %v with error: %w", matches, err)
		}
	}

	i.log.Info("pull image done", "image", image)
	return nil
}

// Burn a image pulling a tarball and unpack to a specific directory
func (i *Image) Burn(prefix, image, source string) error {
	i.log.Info("burn image", "image", image)
	begin := time.Now()

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("%s: failed to open archive %w", source, err)

	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat %s %w", source, err)
	}

	if !strings.HasSuffix(image, "lz4") {
		return fmt.Errorf("unsupported image compression format of image:%s", image)
	}

	lz4Reader := lz4.NewReader(file)
	i.log.Info("lz4", "size", lz4Reader.Size())
	creader := io.NopCloser(lz4Reader)
	// wild guess for lz4 compression ratio
	// lz4 is a stream format and therefore the
	// final size cannot be calculated upfront
	csize := stat.Size() * 2
	defer creader.Close()

	bar := pb.New64(csize)
	bar.Set(pb.Bytes, true)
	bar.Start()
	bar.SetWidth(80)

	reader := bar.NewProxyReader(creader)

	err = archiver.Tar.Read(reader, prefix)
	if err != nil {
		return fmt.Errorf("unable to burn image %s %w", source, err)
	}

	bar.Finish()

	err = os.Remove(source)
	if err != nil {
		i.log.Warn("burn image unable to remove image source", "error", err)
	}

	i.log.Info("burn took", "duration", time.Since(begin))
	return nil
}

// checkHash check the sha512 or md5 signature of file with the sha512sum or md5sum given in the file.
// the content of the file must be in the form:
// <sha512sum | md5sum> filename
// this is the same format as create by the "sha512 | md5sum" unix command
func (i *Image) checkHash(file, hashfile, hashType string) (bool, error) {
	hashfileContent, err := os.ReadFile(hashfile)
	if err != nil {
		return false, fmt.Errorf("unable to read hash file %s %w", hashfile, err)
	}
	expectedHash := strings.Split(string(hashfileContent), " ")[0]

	f, err := os.Open(file)
	if err != nil {
		return false, fmt.Errorf("unable to read file: %s %w", file, err)
	}
	defer f.Close()

	var h hash.Hash
	switch hashType {
	case HashSHA512:
		h = sha512.New()
	case HashMD5:
		h = md5.New()
	default:
		return false, fmt.Errorf("unsupported hash type: %s", hashType)
	}

	if _, err := io.Copy(h, f); err != nil {
		return false, fmt.Errorf("unable to calculate %s of file: %s %w", hashType, file, err)
	}
	sourceHash := fmt.Sprintf("%x", h.Sum(nil))
	i.log.Info("check hash", "source hash", sourceHash, "expected hash", expectedHash)
	if sourceHash != expectedHash {
		return false, fmt.Errorf("source %s:%s expected %s:%s", hashType, sourceHash, hashType, expectedHash)
	}
	return true, nil
}

// downloadFile will download from a source url to a local file dest.
// It's efficient because it will write as it downloads
// and not load the whole file into memory.
func (i *Image) download(source, dest string) error {
	i.log.Info("download", "from", source, "to", dest)
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("unable to create destination %s %w", dest, err)
	}
	defer out.Close()

	// Get the data
	//nolint:gosec,noctx
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("download of %s did not work, statuscode was: %d", source, resp.StatusCode)
	}

	fileSize := resp.ContentLength

	bar := pb.New64(fileSize)
	bar.Set(pb.Bytes, true)
	bar.SetWidth(80)
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
