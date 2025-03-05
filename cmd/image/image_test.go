package image

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"
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

func TestCheckSHA512(t *testing.T) {
	testfile := "/tmp/testsha512"
	testfileSHA512 := "/tmp/testsha512.sha512sum"
	content := []byte("This is testcontent")
	err := os.WriteFile(testfile, content, os.ModePerm) // nolint:gosec
	if err != nil {
		t.Error(err)
	}
	cmd := exec.Command("sha512sum", testfile)
	sha512Content, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}
	sha512, err := os.Create(testfileSHA512)
	if err != nil {
		t.Error(err)
	}
	_, err = sha512.Write(sha512Content)
	if err != nil {
		t.Error(err)
	}
	sha512.Close()
	defer os.Remove(testfile)
	defer os.Remove(testfileSHA512)
	matches, err := NewImage(slog.Default()).checksha512(testfile, testfileSHA512)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("expected sha512 matches, but didn't")
	}
}
