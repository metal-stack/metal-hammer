package main

import (
	"os"
	"time"

	"github.com/metal-stack/go-lldpd/pkg/lldp"
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
