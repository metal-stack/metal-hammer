package image

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	// "github.com/google/go-containerregistry/pkg/name"
	"github.com/stretchr/testify/assert"
)

func TestCheckMD5(t *testing.T) {
	testfile := "/tmp/testmd5"
	testfileMD5 := "/tmp/testmd5.md5"
	content := []byte("This is testcontent")
	err := os.WriteFile(testfile, content, os.ModePerm) // nolint:gosec
	if err != nil {
		t.Error(err)
	}

	cmd := exec.Command("md5sum", testfile)
	md5Content, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}
	md5, err := os.Create(testfileMD5)
	if err != nil {
		t.Error(err)
	}
	_, err = md5.Write(md5Content)
	if err != nil {
		t.Error(err)
	}
	md5.Close()

	defer os.Remove(testfile)
	defer os.Remove(testfileMD5)

	matches, err := NewImage(slog.Default()).checkMD5(testfile, testfileMD5)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("expected md5 matches, but didn't")
	}

}

func TestOciPull(t *testing.T) {
	var (
		assert = assert.New(t)

		ctx                  = context.Background()
		invalidImageRef      = "invalid://"
		imageRef             = "oci://ghcr.io/metal-stack/debian:12-oci-artifact-push" // TODO: change to debian:12 image before merging
		imageRefDoesNotExist = "oci://does/not/exist:tag"
		mountDir             = "/tmp/oci-pull-mount-dir"
		installGoBin         = "install-go"
		anonymousUsername    = ""
		anonymousPassword    = ""

	// TODO: what credentials shall be used here?
	// username = "test-user"
	// password = "test-password"
	)

	t.Run("successful anonymous pull", func(t *testing.T) {
		err := os.Mkdir(mountDir, os.ModePerm)
		if err != nil {
			t.Error(err)
		}
		defer os.RemoveAll(mountDir)

		i := NewImage(slog.Default())
		if err = i.OciPull(ctx, imageRef, mountDir, anonymousUsername, anonymousPassword); err != nil {
			t.Error(err)
		}

		installGoBinFullPath := filepath.Join(mountDir, installGoBin)
		assert.FileExists(installGoBinFullPath)
	})

	// t.Run("successful authenticated pull", func(t *testing.T) {
	// 	err := os.Mkdir(mountDir, os.ModePerm)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	defer os.RemoveAll(mountDir)
	//
	// 	i := NewImage(slog.Default())
	// 	// TODO: what credentials shall be used here?
	// 	if err = i.OciPull(ctx, imageRef, mountDir, username, password); err != nil {
	// 		t.Error(err)
	// 	}
	//
	// 	installGoBinFullPath := filepath.Join(mountDir, installGoBin)
	// 	assert.FileExists(installGoBinFullPath)
	// })

	t.Run("parsing of image refs fails", func(t *testing.T) {
		i := NewImage(slog.Default())
		err := i.OciPull(ctx, invalidImageRef, mountDir, anonymousUsername, anonymousPassword)
		assert.EqualError(err, "parsing image reference: could not parse reference: invalid://")
	})

	t.Run("pulling remote image fails", func(t *testing.T) {
		i := NewImage(slog.Default())
		err := i.OciPull(ctx, imageRefDoesNotExist, mountDir, anonymousUsername, anonymousPassword)
		assert.Error(err)
	})
}
