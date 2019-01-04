package cmd

import (
	"testing"
)

func TestIPToASN(t *testing.T) {
	ipaddress := "10.0.1.2/24"

	asn, err := ipToASN(ipaddress)
	if err != nil {
		t.Errorf("no error expected got:%v", err)
	}

	if asn != 4200000258 {
		t.Errorf("expected 4200000258 got: %d", asn)
	}
}
