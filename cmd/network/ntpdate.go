package network

import (
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

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

func getTime(log *zap.SugaredLogger, servers []string) (t time.Time, err error) {
	for _, s := range servers {
		log.Debugw("ntpdate", "getting time from", s)
		if t, err = ntp.Time(s); err == nil {
			// Right now we return on the first valid time.
			// We can implement better heuristics here.
			log.Debugw("ntpdate", "got time", t)
			return t, nil
		}
	}
	err = fmt.Errorf("unable to get any time from servers %v", servers)
	return
}

// NtpDate set the system time to the time comming from a ntp source
func NtpDate(log *zap.SugaredLogger) {
	t, err := getTime(log, ntpServers)
	if err != nil {
		log.Errorw("ntpdate", "unable to get time", err)
	}

	tv := syscall.NsecToTimeval(t.UnixNano())
	if err = syscall.Settimeofday(&tv); err != nil {
		log.Errorw("ntpdate", "unable to set system time", err)
	}
}
