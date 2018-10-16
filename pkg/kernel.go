package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	cmdline = "/proc/cmdline"
)

// ParseCmdline will put each key=value pair from /proc/cmdline into a map.
func ParseCmdline() (map[string]string, error) {
	cmdline, err := ioutil.ReadFile(cmdline)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %v", cmdline, err)
	}

	cmdLineValues := strings.Split(string(cmdline), " ")
	envmap := make(map[string]string)
	for _, v := range cmdLineValues {
		keyValue := strings.Split(v, "=")
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			envmap[key] = value
		}
	}
	return envmap, nil
}
