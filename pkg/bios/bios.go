package bios

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/inconshreveable/log15"
)

const (
	biosVersion = "/sys/class/dmi/id/bios_version"
	biosVendor  = "/sys/class/dmi/id/bios_vendor"
	biosDate    = "/sys/class/dmi/id/bios_date"
)

// BIOS information of this machine
type BIOS struct {
	Version string
	Vendor  string
	Date    string
}

// Bios read bios informations
func Bios() *BIOS {
	return &BIOS{
		Version: read(biosVersion),
		Vendor:  read(biosVendor),
		Date:    read(biosDate),
	}
}

func (b *BIOS) String() string {
	return "version:" + b.Version + " vendor:" + b.Vendor + " date:" + b.Date
}

func read(file string) string {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Error("error reading", "file", file, "error", err)
			return ""
		}
		return strings.TrimSpace(string(content))
	}
	return ""
}
