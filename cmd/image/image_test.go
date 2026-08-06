package image

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"

	//nolint:gosec
	"crypto/md5"
	"crypto/sha512"
)

func TestCheckMD5(t *testing.T) {
	testfile := "/tmp/testmd5"
	testfileMD5 := "/tmp/testmd5.md5File"
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
	md5File, err := os.Create(testfileMD5)
	if err != nil {
		t.Error(err)
	}
	_, err = md5File.Write(md5Content)
	if err != nil {
		t.Error(err)
	}
	md5File.Close()

	defer os.Remove(testfile)
	defer os.Remove(testfileMD5)

	matches, err := NewImage(slog.Default()).checkHash(testfile, testfileMD5, md5.New)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("expected md5File matches, but didn't")
	}

}

func TestCheckSHA512(t *testing.T) {
	testfile := "/tmp/testsha512"
	testfileSHA512 := "/tmp/testsha512.sha512sum"
	content := []byte("This is testcontent")
	err := os.WriteFile(testfile, content, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	cmd := exec.Command("sha512sum", testfile)
	sha512Content, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}
	sha512File, err := os.Create(testfileSHA512)
	if err != nil {
		t.Error(err)
	}
	_, err = sha512File.Write(sha512Content)
	if err != nil {
		t.Error(err)
	}
	sha512File.Close()
	defer os.Remove(testfile)
	defer os.Remove(testfileSHA512)
	matches, err := NewImage(slog.Default()).checkHash(testfile, testfileSHA512, sha512.New)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("expected sha512File matches, but didn't")
	}
}
