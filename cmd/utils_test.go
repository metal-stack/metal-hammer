package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestCheckMD5(t *testing.T) {
	testfile := "/tmp/testmd5"
	testfileMD5 := "/tmp/testmd5.md5"
	content := []byte("This is testcontent")
	err := ioutil.WriteFile(testfile, content, os.ModePerm)
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
	md5.Write(md5Content)
	md5.Close()

	defer os.Remove(testfile)
	defer os.Remove(testfileMD5)

	matches, err := checkMD5(testfile, testfileMD5)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("expected md5 matches, but didn't")
	}

}
