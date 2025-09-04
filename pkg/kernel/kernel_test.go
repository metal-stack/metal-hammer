package kernel

import (
	"fmt"
	"os"
	"testing"
)

func TestParseCmdLine(t *testing.T) {
	cmdline = "/tmp/testcmdline"
	defer os.Remove(cmdline)

	tt := []struct {
		commandline string
		result      [][]string
	}{
		{"quiet", [][]string{}},
		{"console=ttyS0", [][]string{{"console", "ttyS0"}}},
		{"console=tty root=/dev/sda1", [][]string{{"console", "tty"}, {"root", "/dev/sda1"}}},
	}

	for _, tc := range tt {

		err := writeCmdline(tc.commandline)
		if err != nil {
			t.Error(err)
		}

		envpairs, err := ParseCmdline()
		if err != nil {
			t.Error(err)
		}
		for key := range tc.result {
			switch envpairs[key][0] {
			case "console":
				if envpairs[key][1] != tc.result[key][1] {
					t.Errorf("expected %s but got %s", tc.result[key], envpairs[key][1])
				}
			case "root":
				if envpairs[key][1] != tc.result[key][1] {
					t.Errorf("expected %s but got %s", tc.result[key], envpairs[key][1])
				}
			default:
			}
		}
	}
}

func TestFirmware(t *testing.T) {
	sysfirmware = "/tmp/testefi"
	_, err := os.OpenFile(sysfirmware, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(sysfirmware)

	firmware := Firmware()
	if firmware != "efi" {
		t.Error("expected efi firmware but didn't get")
	}

	sysfirmware = "/tmp/testbios"
	firmware = Firmware()
	if firmware != "bios" {
		t.Error("expected bios firmware but didn't get")
	}
}

func writeCmdline(content string) error {
	err := os.WriteFile(cmdline, []byte(content), os.ModePerm) // nolint:gosec
	if err != nil {
		return fmt.Errorf("unable to write test cmdline")
	}
	return nil
}
