package service

import (
	"fmt"
	"net"
)

func GetPrefix(inIP string, inMask string) (string, error) {
	ipv4 := net.ParseIP(inIP)
	if ipv4 == nil {
		return "", fmt.Errorf("invalid IP address %s", inIP)
	}
	mask := net.ParseIP(inMask)
	if mask == nil {
		return "", fmt.Errorf("invalid IP mask: %s", inMask)
	}
	for i := range ipv4 {
		ipv4[i] &= mask[i]
	}
	return ipv4.String(), nil
}
