package image

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/mholt/archiver"
	lz4 "github.com/pierrec/lz4/v4"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	//nolint:gosec
	"crypto/md5"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Image struct {
	log *slog.Logger
}

func NewImage(log *slog.Logger) *Image {
	return &Image{log: log}
}

func (i *Image) OciPull(ctx context.Context, imageRef, mountDir, username, password string) error {
	// Parse the image reference (e.g., docker.io/library/alpine:latest)
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("parsing image reference: %w", err)
	}

	// Choose authentication method
	var auth = authn.Anonymous
	if username != "" || password != "" {
		auth = &authn.Basic{
			Username: username,
			Password: password,
		}
	}

	i.log.Info("pull oci image", "image", imageRef)
	img, err := remote.Image(ref, remote.WithAuth(auth))
	if err != nil {
		return fmt.Errorf("fetching remote image: %w", err)
	}

	// Flatten layers and create a tar stream
	rc := mutate.Extract(img)
	defer rc.Close()

	i.log.Info(fmt.Sprintf("untar oci image into %s", mountDir), "image", imageRef)
	if err := i.untar(rc, mountDir); err != nil {
		return fmt.Errorf("extracting tar: %w", err)
	}

	i.log.Info("pull oci image done", "image", imageRef)
	return nil
}

// Pull an image from s3
func (i *Image) Pull(image, destination string) error {
	i.log.Info("pull image", "image", image)
	md5destination := destination + ".md5"
	md5file := image + ".md5"
	err := i.download(image, destination)
	if err != nil {
		return fmt.Errorf("unable to pull image %s %w", image, err)
	}
	err = i.download(md5file, md5destination)
	defer os.Remove(md5destination)
	if err != nil {
		return fmt.Errorf("unable to pull md5 %s %w", md5file, err)
	}
	i.log.Info("check md5")
	matches, err := i.checkMD5(destination, md5destination)
	if err != nil || !matches {
		return fmt.Errorf("md5sum mismatch")
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

// checkMD5 check the md5 signature of file with the md5sum given in the md5file.
// the content of the md5file must be in the form:
// <md5sum> filename
// this is the same format as create by the "md5sum" unix command
func (i *Image) checkMD5(file, md5file string) (bool, error) {
	md5fileContent, err := os.ReadFile(md5file)
	if err != nil {
		return false, fmt.Errorf("unable to read md5sum file %s %w", md5file, err)
	}
	expectedMD5 := strings.Split(string(md5fileContent), " ")[0]

	f, err := os.Open(file)
	if err != nil {
		return false, fmt.Errorf("unable to read file: %s %w", file, err)
	}
	defer f.Close()

	//nolint:gosec
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, fmt.Errorf("unable to calculate md5sum of file: %s %w", file, err)
	}
	sourceMD5 := fmt.Sprintf("%x", h.Sum(nil))
	i.log.Info("check md5", "source md5", sourceMD5, "expected md5", expectedMD5)
	if sourceMD5 != expectedMD5 {
		return false, fmt.Errorf("source md5:%s expected md5:%s", sourceMD5, expectedMD5)
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

func (i *Image) untar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}

		target := filepath.Join(dest, hdr.Name)

		fmt.Printf("extracting:%s\n", target)

		if strings.HasSuffix(target, ".log") {
			fmt.Printf("skipping:%s\n", target)
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return fmt.Errorf("creating dir: %w", err)
			}
			if err := os.Lchown(target, hdr.Uid, hdr.Gid); err != nil && !errors.Is(err, os.ErrPermission) {
				return fmt.Errorf("chown dir: %w", err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("creating parent dir: %w", err)
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return fmt.Errorf("creating file: %w", err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("copying file: %w", err)
			}
			f.Close()

			if err := os.Chmod(target, os.FileMode(hdr.Mode)); err != nil {
				return fmt.Errorf("chmod: %w", err)
			}
			if err := os.Lchown(target, hdr.Uid, hdr.Gid); err != nil && !errors.Is(err, os.ErrPermission) {
				return fmt.Errorf("chown file: %w", err)
			}

		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("creating parent dir: %w", err)
			}
			if err := os.Symlink(hdr.Linkname, target); err != nil {
				return fmt.Errorf("creating symlink: %w", err)
			}
			if err := os.Lchown(target, hdr.Uid, hdr.Gid); err != nil && !errors.Is(err, os.ErrPermission) {
				return fmt.Errorf("chown symlink: %w", err)
			}

		default:
			// skip unsupported or special files
		}
	}
	return nil
}
