package hardware

import (
	"bytes"
	"io/ioutil"
	"net"
	"path/filepath"
	"sort"
	"strings"
)

// 前提: bonding的mac是子网卡的mac之一(目前观察到是这样, 但未知bonding是否可自定义mac, 下面代码是基于bongding可能有自定义mac来实现), 则遍历网络接口再去虚拟接口去重即可得到RealMACs.
func RealMACs() ([]string, error) {
	// 只取网络接口, 因此避免了类似＂ifconfig/python的netifaces.interfaces()将ip alias显示为新网络接口＂
	interfaces, err := net.Interfaces() // from /sys/class/net
	if err != nil {
		return nil, err
	}

	macs := make([]string, 0, len(interfaces))
	m := make(map[string]string)

	for _, v := range interfaces {
		if FileIsDir(filepath.Join("/sys/devices/virtual/net", v.Name)) { // skip vitrual nic(includ bond,lo)
			continue
		}

		m[v.Name] = v.HardwareAddr.String()
	}

	// add bond slaves
	for _, v := range interfaces {
		if !FileIsExist(filepath.Join("/proc/net/bonding", v.Name)) {
			continue
		}

		delete(m, v.Name) // bond'mac 默认是参与bond的之一网卡的mac, 这里防止bond mac与所有参与bond的mac都不一样.

		slaves, err := ParseBondNIC(v.Name)
		if err != nil {
			return nil, err
		}

		for bk, bv := range slaves {
			m[bk] = bv
		}
	}

	for _, v := range m {
		macs = append(macs, v)
	}

	sort.Strings(macs)

	return macs, nil
}

// bond group use the same mac, so need parse
func ParseBondNIC(name string) (map[string]string, error) {
	m := make(map[string]string)

	slavesStr, err := ioutil.ReadFile(filepath.Join("/sys/devices/virtual/net", name, "bonding/slaves"))
	if err != nil {
		return nil, err
	}

	slaves := strings.Fields(strings.TrimSpace(string(slavesStr)))

	var tmp []byte
	for _, v := range slaves {
		tmp, err = ioutil.ReadFile(filepath.Join("/sys/devices/virtual/net", name, "lower_"+v, "bonding_slave/perm_hwaddr"))
		if err != nil {
			return nil, err
		}

		if len(tmp) > 0 {
			m[v] = string(bytes.TrimSpace(tmp))
		}
	}

	return m, nil
}
