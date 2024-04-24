package utilities

import (
	"fmt"
	"net"
)

func ValidateMACAddress(mac string) error {
	_, err := net.ParseMAC(mac)

	if err != nil {
		return fmt.Errorf("invalid MAC address")
	}

	return nil
}

func ValidateIPAddress(ip string) error {
	parsed := net.ParseIP(ip)

	if parsed == nil {
		return fmt.Errorf("invalid IP address")
	}

	if parsed.IsUnspecified() || parsed.To4()[0] == 0x00 {
		return fmt.Errorf("invalid IP address range")
	}

	return nil
}

func ValidateSubnetMask(mask string) error {
	ip := net.ParseIP(mask).To4()
	_, bits := net.IPMask(ip).Size()
	if bits == 0 { // Invalid subnet masks will return 0 bits
		return fmt.Errorf("invalid subnet mask")
	}

	return nil
}
