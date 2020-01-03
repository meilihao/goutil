package hardware

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

func RealMACs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	macs := make([]string, 0, len(interfaces))
	m := make(map[string]string)

	for _, v := range interfaces {
		if FileIsDir(filepath.Join("/sys/devices/virtual/net", v.Name)) { // skip vitrual nic
			continue
		}

		m[v.Name] = v.HardwareAddr.String()
	}

	for _, v := range interfaces {
		if !FileIsExist(filepath.Join("/proc/net/bonding", v.Name)) {
			continue
		}

		for bk, bv := range ParseBondNIC(v.Name) {
			m[bk] = bv
		}
	}

	for _, v := range m {
		macs = append(macs, v)
	}

	sort.Strings(macs)

	return macs
}

// bond group use the same mac, so need parse
func ParseBondNIC(name string) map[string]string {
	m := make(map[string]string)

	slavesStr, err := ioutil.ReadFile(filepath.Join("/sys/devices/virtual/net", name, "bonding/slaves"))
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	slaves := strings.Fields(string(slavesStr))

	var tmp []byte
	for _, v := range slaves {
		tmp, err = ioutil.ReadFile(filepath.Join("/sys/devices/virtual/net", name, "lower_"+v, "bonding_slave/perm_hwaddr"))
		if err != nil {
			log.Error().Err(err)
			return nil
		}

		if len(tmp) > 0 {
			m[v] = string(tmp)

			tmp = tmp[:0]
		}
	}

	return m
}

func FileIsExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FileIsDir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
