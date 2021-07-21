package util

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/jumperzq86/jumper_conn/def"
)

func GetExternalIp() (string, error) {
	rsp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", def.ErrGetExternalIp
	}
	ipbytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(ipbytes))
	return ip, nil
}

func GetInternalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	//
	return "", nil
}

func GetInternalIps() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips, nil
}
