package image

import (
	"archive/tar"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestUntar(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name     string
		fileType byte
		mode     os.FileMode
		content  string
	}{
		{"setuid_file", tar.TypeReg, 04755, "test content"}, // Regular file with setuid bit
		{"setgid_file", tar.TypeReg, 02755, "test content"}, // Regular file with setgid bit
		{"sticky_file", tar.TypeReg, 01755, "test content"}, // Regular file with sticky bit
		{"setgid_dir", tar.TypeDir, 02755, ""},              // Directory with setgid bit
		{"sticky_dir", tar.TypeDir, 01755, ""},              // Directory with sticky bit
		{"symlink", tar.TypeSymlink, 0777, "target"},        // Symbolic link
		// {"shadow_file", tar.TypeReg, 0000, "test content"},  // Simulate special files - FIXME: special files need root permissions
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tarPath := filepath.Join(tempDir, tc.name+".tar")
			createTestTar(t, tarPath, tc.name, tc.fileType, tc.mode, tc.content)

			destDir := filepath.Join(tempDir, "extracted_"+tc.name)
			i := NewImage(slog.Default())

			tarFile, err := os.Open(tarPath)
			require.NoError(t, err)
			defer tarFile.Close()
			err = i.untar(tarFile, destDir)
			require.NoError(t, err)

			extractedPath := filepath.Join(destDir, tc.name)
			verifyExtractedFile(t, extractedPath, tc.fileType, tc.mode, tc.content)
		})
	}

	// TODO: consider character and block devices -> skipped for now due to root permissions
}

func createTestTar(t *testing.T, tarPath, fileName string, fileType byte, mode os.FileMode, content string) {
	tarFile, err := os.Create(tarPath)
	require.NoError(t, err)
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	header := &tar.Header{
		Name:     fileName,
		Typeflag: fileType,
	}
	if fileType == tar.TypeSymlink {
		header.Linkname = content
	} else {
		header.Size = int64(len(content))
		header.Mode = int64(mode)
	}

	err = tw.WriteHeader(header)
	require.NoError(t, err)

	if fileType == tar.TypeReg {
		_, err := tw.Write([]byte(content))
		require.NoError(t, err)
	}
}

func verifyExtractedFile(t *testing.T, path string, fileType byte, mode os.FileMode, content string) {
	info, err := os.Lstat(path)
	require.NoError(t, err)
	require.Equal(t, info.Mode().Perm(), mode.Perm())

	switch fileType {
	case tar.TypeReg:
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, string(data), content)
	case tar.TypeSymlink:
		target, err := os.Readlink(path)
		require.NoError(t, err)
		require.Equal(t, target, content)
	case tar.TypeDir:
		require.Equal(t, info.IsDir(), true)
	}
}
