package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/google/uuid"
	log "github.com/inconshreveable/log15"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// downloadFile will download from a source url to a local file dest.
// It's efficient because it will write as it downloads
// and not load the whole file into memory.
func download(source, dest string) error {
	log.Info("download", "from", source, "to", dest)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("download of %s did not work, statuscode was: %d", source, resp.StatusCode)
	}

	fileSize := resp.ContentLength

	bar := pb.New64(fileSize).SetUnits(pb.U_BYTES)
	bar.SetWidth(80)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	reader := bar.NewProxyReader(resp.Body)
	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	return nil
}

// checkMD5 check the md5 signature of file with the md5sum given in the md5file.
// the content of the md5file must be in the form:
// <md5sum> filename
// this is the same format as create by the "md5sum" unix command
func checkMD5(file, md5file string) (bool, error) {
	md5fileContent, err := ioutil.ReadFile(md5file)
	if err != nil {
		return false, fmt.Errorf("unable to read md5sum file: %v", err)
	}
	expectedMD5 := strings.Split(string(md5fileContent), " ")[0]

	f, err := os.Open(file)
	if err != nil {
		return false, fmt.Errorf("unable to read file: %v", err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, fmt.Errorf("unable to calculate md5sum of file: %v", err)
	}
	sourceMD5 := fmt.Sprintf("%x", h.Sum(nil))
	log.Info("checkMD5", "source md5", sourceMD5, "expected md5", expectedMD5)
	if sourceMD5 != expectedMD5 {
		return false, fmt.Errorf("source md5:%s expected md5:%s", sourceMD5, expectedMD5)
	}
	return true, nil
}

// small helper to execute a command, redirect stdout/stderr.
func executeCommand(name string, arg ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("unable to locate program:%s in path info:%v", name, err)
	}
	cmd := exec.Command(path, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

// we start to calculate ASNs for devices with the first ASN in the 32bit ASN range and
// add the last 2 octets of the ip of the device to achieve unique ASNs per vrf
func ipToASN(ipaddress string) (int64, error) {
	const asnbase = 4200000000

	ip, _, err := net.ParseCIDR(ipaddress)
	if err != nil {
		return int64(-1), fmt.Errorf("unable to parse ip %v", err)
	}

	asn := asnbase + int64(ip[14])*256 + int64(ip[15])
	return asn, nil
}

// save the content of kernel ringbuffer to /var/log/syslog
// by calling the appropriate syscall.
func createSyslog() error {
	const SyslogActionReadAll = 3
	level := uintptr(SyslogActionReadAll)

	b := make([]byte, 256*1024)
	amt, _, err := syscall.Syscall(syscall.SYS_SYSLOG, level, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if err != 0 {
		return err
	}

	return ioutil.WriteFile("/var/log/syslog", b[:amt], 0666)
}

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
			productSerialBytes, err := uuid.FromBytes([]byte(fmt.Sprintf("%16s", string(productSerial))))
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

func InternalIP() string {
	var ip net.IP
	interfaces := []string{"eth0", "eth1", "eth2", "eth3", "eth4", "eth5", "eth6", "eth7", "eth8", "eth9"}
	for _, eth := range interfaces {
		itf, _ := net.InterfaceByName(eth)
		item, _ := itf.Addrs()
		for _, addr := range item {
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						ip = v.IP
					}
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	}
	return ""
}
