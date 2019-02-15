package cmd

import (
	"github.com/pkg/errors"
	"io"
	"net"
	"os"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.Errorf("%s is not a regular file", src)
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

// we start to calculate ASNs for machines with the first ASN in the 32bit ASN range and
// add the last 2 octets of the ip of the machine to achieve unique ASNs per vrf
func ipToASN(ipaddress string) (int64, error) {
	const asnbase = 4200000000

	ip, _, err := net.ParseCIDR(ipaddress)
	if err != nil {
		return int64(-1), errors.Wrapf(err, "unable to parse ip %s", ipaddress)
	}

	asn := asnbase + int64(ip[14])*256 + int64(ip[15])
	return asn, nil
}
