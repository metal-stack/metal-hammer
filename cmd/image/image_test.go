package image

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
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

		ctx          = context.Background()
		mountDir     = "/tmp/oci-pull-mount-dir"
		extractedBin = "a"

		anonymousUsername = ""
		anonymousPassword = ""
	)

	t.Run("successful anonymous pull", func(t *testing.T) {
		regIP, regPort, err := startRegistry(nil, nil, nil)
		require.NoError(t, err)
		registry := fmt.Sprintf("%s:%d", regIP, regPort)

		imageRef := fmt.Sprintf("%s/library/image", registry)
		err = createImage(imageRef, "", "")
		require.NoError(t, err)

		err = os.MkdirAll(mountDir, 0777)
		require.NoError(t, err)
		defer os.RemoveAll(mountDir)

		i := NewImage(slog.Default())
		if err = i.OciPull(ctx, imageRef, mountDir, anonymousUsername, anonymousPassword); err != nil {
			t.Error(err)
		}

		extractedBinFullPath := filepath.Join(mountDir, extractedBin)
		assert.FileExists(extractedBinFullPath)
	})

	t.Run("successful authenticated pull", func(t *testing.T) {
		var (
			username = "test-user"
			password = "test-password"
		)

		f, err := os.CreateTemp("", "htpasswd")
		require.NoError(t, err)
		defer func() {
			_ = os.Remove(f.Name())
		}()

		err = htpasswd.SetPassword(f.Name(), username, password, htpasswd.HashBCrypt)
		require.NoError(t, err)

		env := map[string]string{
			"REGISTRY_AUTH":                "htpasswd",
			"REGISTRY_AUTH_HTPASSWD_REALM": "registry-login",
			"REGISTRY_AUTH_HTPASSWD_PATH":  "/htpasswd",
		}
		regIP, regPort, err := startRegistry(env, pointer.Pointer(f.Name()), pointer.Pointer("/htpasswd"))
		require.NoError(t, err)
		registry := fmt.Sprintf("%s:%d", regIP, regPort)

		imageRefBehindAuth := fmt.Sprintf("%s/library/image", registry)
		err = createImage(imageRefBehindAuth, username, password)
		require.NoError(t, err)

		err = os.MkdirAll(mountDir, 0777)
		require.NoError(t, err)
		defer os.RemoveAll(mountDir)

		i := NewImage(slog.Default())
		if err = i.OciPull(ctx, imageRefBehindAuth, mountDir, username, password); err != nil {
			t.Error(err)
		}

		extractedBinFullPath := filepath.Join(mountDir, extractedBin)
		assert.FileExists(extractedBinFullPath)
	})

	t.Run("parsing of image refs fails", func(t *testing.T) {
		invalidImageRef := "invalid://"
		i := NewImage(slog.Default())
		err := i.OciPull(ctx, invalidImageRef, mountDir, anonymousUsername, anonymousPassword)
		assert.EqualError(err, "parsing image reference: could not parse reference: invalid://")
	})

	t.Run("pulling remote image fails", func(t *testing.T) {
		imageRefDoesNotExist := "oci://does/not/exist:tag"
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

func createImage(imageName, username, password string, tags ...string) error {
	// ensure every image has distinct content
	buf := make([]byte, 128)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}
	img, err := crane.Image(map[string][]byte{"a": buf})
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
	for _, tag := range tags {
		err := crane.Push(img, imageName+":"+tag, crane.WithAuth(auth))
		if err != nil {
			return err
		}
	}

	return nil
}
