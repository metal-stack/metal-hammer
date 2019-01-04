package uuid

import (
	"fmt"
	guuid "github.com/google/uuid"
	log "github.com/inconshreveable/log15"
	"io/ioutil"
	"os"
	"strings"
)

const dmiUUID = "/sys/class/dmi/id/product_uuid"
const dmiSerial = "/sys/class/dmi/id/product_serial"

// DeviceUUID calculates a unique uuid for this (hardware) device
func DeviceUUID() string {
	if _, err := os.Stat(dmiUUID); !os.IsNotExist(err) {
		productUUID, err := ioutil.ReadFile(dmiUUID)
		if err != nil {
			log.Error("error getting product_uuid", "error", err)
		} else {
			log.Info("create UUID from", "source", dmiUUID)
			return strings.TrimSpace(string(productUUID))
		}
	}

	if _, err := os.Stat(dmiSerial); !os.IsNotExist(err) {
		productSerial, err := ioutil.ReadFile(dmiSerial)
		if err != nil {
			log.Error("error getting product_serial", "error", err)
		} else {
			productSerialBytes, err := guuid.FromBytes([]byte(fmt.Sprintf("%16s", string(productSerial))))
			if err != nil {
				log.Error("error getting converting product_serial to uuid", "error", err)
			} else {
				log.Info("create UUID from", "source", dmiSerial)
				return strings.TrimSpace(productSerialBytes.String())
			}
		}
	}
	log.Error("no valid UUID found", "return uuid", "00000000-0000-0000-0000-000000000000")
	return "00000000-0000-0000-0000-000000000000"
}
