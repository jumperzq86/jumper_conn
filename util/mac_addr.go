package util

import (
	"net"

	"github.com/jumperzq86/jumper_conn/def"
)

func GetMacAddr() (string, error) {
	inters, err := net.Interfaces()
	if err != nil {
		return "", def.ErrGetMacAddr
	}

	for _, inter := range inters {
		macAddr := inter.HardwareAddr.String()
		if len(macAddr) != 0 {
			return macAddr, nil
		}
	}

	return "", nil
}

func GetMacAddrs() ([]string, error) {
	inters, err := net.Interfaces()
	if err != nil {
		return nil, def.ErrGetMacAddr
	}

	macAddrs := make([]string, 0)
	for _, inter := range inters {
		macAddr := inter.HardwareAddr.String()
		if len(macAddr) != 0 {
			macAddrs = append(macAddrs, macAddr)
		}
	}

	return macAddrs, nil
}
