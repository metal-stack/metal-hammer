package image

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/foomo/htpasswd"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

		anonymousUsername = ""
		anonymousPassword = ""
		authUsername      = "test-user"
		authPassword      = "test-password"
	)

	f, err := os.CreateTemp("", "htpasswd")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(f.Name())
	}()

	err = htpasswd.SetPassword(f.Name(), authUsername, authPassword, htpasswd.HashBCrypt)
	require.NoError(t, err)

	env := map[string]string{
		"REGISTRY_AUTH":                "htpasswd",
		"REGISTRY_AUTH_HTPASSWD_REALM": "registry-login",
		"REGISTRY_AUTH_HTPASSWD_PATH":  "/htpasswd",
	}
	regIP, regPort, err := startRegistry(env, pointer.Pointer(f.Name()), pointer.Pointer("/htpasswd"))
	require.NoError(t, err)
	registry := fmt.Sprintf("%s:%d", regIP, regPort)

	imageRefBehindAuth := fmt.Sprintf("%s/library/debian", registry)
	trimmedImageRef := strings.TrimPrefix(imageRef, "oci://")
	err = fetchImageFromRemote(imageRefBehindAuth, trimmedImageRef, authUsername, authPassword)
	require.NoError(t, err)

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

	t.Run("successful authenticated pull", func(t *testing.T) {
		err := os.Mkdir(mountDir, os.ModePerm)
		if err != nil {
			t.Error(err)
		}
		defer os.RemoveAll(mountDir)

		i := NewImage(slog.Default())
		if err = i.OciPull(ctx, imageRefBehindAuth, mountDir, authUsername, authPassword); err != nil {
			t.Error(err)
		}

		installGoBinFullPath := filepath.Join(mountDir, installGoBin)
		assert.FileExists(installGoBinFullPath)
	})

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

// HELPER FUNCTIONS
func startRegistry(env map[string]string, src, dst *string) (string, int, error) {
	ctx := context.Background()
	var (
		c   testcontainers.Container
		err error
	)

	req := testcontainers.ContainerRequest{
		Image:        "registry:3",
		ExposedPorts: []string{"5000/tcp"},
		Env:          env,
		WaitingFor: wait.ForAll(
			wait.ForLog("listening on"),
			wait.ForListeningPort("5000/tcp"),
		),
	}
	c, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", 0, err
	}
	if src != nil && dst != nil {
		err = c.CopyFileToContainer(ctx, *src, *dst, 0o777)
		if err != nil {
			return "", 0, err
		}
	}

	ip, err := c.Host(ctx)
	if err != nil {
		return ip, 0, err
	}
	port, err := c.MappedPort(ctx, "5000")
	if err != nil {
		return ip, port.Int(), err
	}

	return ip, port.Int(), nil
}

func fetchImageFromRemote(imageName, remoteImageName, username, password string) error {
	img, err := crane.Pull(remoteImageName)
	if err != nil {
		return err
	}

	var auth = authn.Anonymous
	if username != "" || password != "" {
		auth = &authn.Basic{
			Username: username,
			Password: password,
		}
	}
	err = crane.Push(img, imageName, crane.WithAuth(auth))
	if err != nil {
		return err
	}

	return nil
}
