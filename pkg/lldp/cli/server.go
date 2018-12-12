package main

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/lldp"
	"os"
	"time"
)

func main() {
	iface := os.Args[1]

	lldpd, err := lldp.NewDaemon("metal-hammer", "waiting for installation", iface, 2*time.Second)
	if err != nil {
		panic(err)
	}
	lldpd.Start()
	select {}
}
