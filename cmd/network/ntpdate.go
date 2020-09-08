package network

import (
	"fmt"
	"syscall"
	"time"

	log "github.com/inconshreveable/log15"

	"github.com/beevik/ntp"
)

var (
	ntpServers = []string{
		"0.de.pool.ntp.org",
		"1.de.pool.ntp.org",
		"2.de.pool.ntp.org",
		"3.de.pool.ntp.org",
		"time.google.com",
	}
)

func getTime(servers []string) (t time.Time, err error) {
	for _, s := range servers {
		log.Debug("ntpdate", "getting time from", s)
		if t, err = ntp.Time(s); err == nil {
			// Right now we return on the first valid time.
			// We can implement better heuristics here.
			log.Debug("ntpdate", "got time", t)
			return t, nil
		}
	}
	err = fmt.Errorf("unable to get any time from servers %v", servers)
	return
}

// NtpDate set the system time to the time comming from a ntp source
func NtpDate() {
	t, err := getTime(ntpServers)
	if err != nil {
		log.Error("ntpdate", "unable to get time", err)
	}

	tv := syscall.NsecToTimeval(t.UnixNano())
	if err = syscall.Settimeofday(&tv); err != nil {
		log.Error("ntpdate", "unable to set system time", err)
	}
}
