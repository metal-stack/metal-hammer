package main

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// FetchAndExtract pulls a container image and extracts its contents into mountDir.
func FetchAndExtract(ctx context.Context, imageRef, mountDir, username, password string) error {
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

	// Pull image
	fmt.Printf("Pulling image %s...\n", imageRef)
	img, err := remote.Image(ref, remote.WithAuth(auth))
	if err != nil {
		return fmt.Errorf("fetching remote image: %w", err)
	}

	// Flatten layers and create a tar stream
	rc := mutate.Extract(img)
	defer rc.Close()

	// Untar into the mountDir
	fmt.Printf("Extracting image into %s...\n", mountDir)
	if err := untar(rc, mountDir); err != nil {
		return fmt.Errorf("extracting tar: %w", err)
	}

	fmt.Println("✅ Extraction complete!")
	return nil
}

func untar(r io.Reader, dest string) error {
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

func main() {
	ctx := context.Background()

	imageRef := "ghcr.io/metal-stack/debian:12-update-kernels"
	mountDir := "/home/stefan/tmp/metal-debian"

	// Leave username/password empty for public registries
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

	if err := FetchAndExtract(ctx, imageRef, mountDir, username, password); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
