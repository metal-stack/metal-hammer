package kernel

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseCmdLine(t *testing.T) {
	cmdline = "/tmp/testcmdline"
	defer os.Remove(cmdline)

	tt := []struct {
		commandline string
		result      map[string]string
	}{
		{"quiet", map[string]string{}},
		{"console=ttyS0", map[string]string{"console": "ttyS0"}},
		{"console=tty root=/dev/sda1", map[string]string{"console": "tty", "root": "/dev/sda1"}},
	}

	for _, tc := range tt {

		err := writeCmdline(tc.commandline)
		if err != nil {
			t.Error(err)
		}

		envmap, err := ParseCmdline()
		if err != nil {
			t.Error(err)
		}
		for key := range tc.result {
			val, ok := envmap[key]
			if !ok {
				t.Error("key not found")
			}
			if val != tc.result[key] {
				t.Errorf("expected %s but got %s", tc.result[key], val)
			}
		}
	}
}

func TestFirmware(t *testing.T) {
	sysfirmware = "/tmp/testefi"
	os.OpenFile(sysfirmware, os.O_RDONLY|os.O_CREATE, 0666)
	defer os.Remove(sysfirmware)

	firmware := Firmware()
	if firmware != "efi" {
		t.Error("expected efi firmware but didnt get")
	}

	sysfirmware = "/tmp/testbios"
	firmware = Firmware()
	if firmware != "bios" {
		t.Error("expected bios firmware but didnt get")
	}
}

func writeCmdline(content string) error {
	err := ioutil.WriteFile(cmdline, []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write test cmdline")
	}
	return nil
}
