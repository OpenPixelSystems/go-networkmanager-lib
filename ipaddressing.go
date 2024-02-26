package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/Wifx/gonetworkmanager/v2"
)

const (
	ip_addr           = "192.168.3.43" // IP address you want to set
	defgateway        = "192.168.3.1"  // Gateway you want to set
	yourInterfaceName = "eth0"         // Interface name you want to use
	prefix_nr         = 24             // Prefix number
	ip_mode           = "manual"       // Mode you want to set

	connectionSection      = "connection"
	connectionSectionID    = "id"
	ip4Section             = "ipv4"
	ip4SectionAddresses    = "addresses"
	ip4SectionAddress      = "address"
	ip4SectionPrefix       = "prefix"
	ip4SectionMethod       = "method"
	ip4SectionNeverDefault = "never-default"
	ip6Section             = "ipv6"
	ip6SectionMethod       = "method"
	connectionID           = "Wired connection 1"
)

func ipAddrToDecimal(ipAddr string) uint32 {
	ip := net.ParseIP(ipAddr)
	decimal := binary.BigEndian.Uint32(ip.To4())

	// Reverse the endianness
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, decimal)
	reversed := binary.BigEndian.Uint32(bytes)
	return reversed
}

func get_original_interface_setting() {
	log.Print("getting original interface settings")
	// Create a new instance of gonetworkmanager
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Get the list of all network devices
	devices, err := nm.GetDevices()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Find the device with the specified interface name
	var selectedDevice gonetworkmanager.Device
	for _, device := range devices {
		interfaceName, _ := device.GetPropertyInterface()
		if interfaceName == yourInterfaceName { // Change to the interface name you want to use
			selectedDevice = device
			break
		}
	}

	if selectedDevice == nil {
		fmt.Println("Selected device not found")
		os.Exit(1)
	}

	// Get the original IPv4 settings of the selected device
	originalIPv4Settings, err := selectedDevice.GetPropertyIP4Config()
	if err != nil {
		fmt.Println("Failed to get original IPv4 settings:", err)
		os.Exit(1)
	}

	// Extract IPv4 addresses
	originalIPv4Addresses, err := originalIPv4Settings.GetPropertyAddressData()
	if err != nil {
		fmt.Println("Failed to get original IPv4 addresses:", err)
		os.Exit(1)
	}

	// Extract the first IPv4 address and its prefix (assuming there's at least one address)
	var originalIPv4Address string
	var originalIPv4Prefix uint
	if len(originalIPv4Addresses) > 0 {
		originalIPv4Address = originalIPv4Addresses[0].Address
		originalIPv4Prefix = uint(originalIPv4Addresses[0].Prefix)
	} else {
		fmt.Println("No IPv4 addresses found in the original settings")
		os.Exit(1)
	}

	// Get the original gateway
	originalGateway, err := originalIPv4Settings.GetPropertyGateway()
	if err != nil {
		fmt.Println("Failed to get original gateway:", err)
		os.Exit(1)
	}

	fmt.Printf("Original IPv4 address: %s/%d\n", originalIPv4Address, originalIPv4Prefix)
	fmt.Printf("Original gateway: %s\n\n", originalGateway)
	return
}

func setIPAddr() {
	log.Print("setting ip address")

	// Create a new instance of gonetworkmanager.Settings
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		fmt.Print("could not get new settings")
		os.Exit(1)
	}

	// Get the list of all connections
	currentConnections, err := settings.ListConnections()
	if err != nil {
		fmt.Print("could not get settings connections list")
		os.Exit(1)
	}

	// Loop through the connections and find the one with the specified ID
	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			fmt.Print("could not get settings of connection")
			os.Exit(1)
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			addressData := make([]map[string]interface{}, 1)
			addressData[0] = make(map[string]interface{})
			addressData[0][ip4SectionPrefix] = 24
			addressData[0][ip4SectionAddress] = ip_addr

			// order defined by network manager
			addresses := make([]uint32, 3)
			// IP addr
			addresses[0] = ipAddrToDecimal(ip_addr)
			addresses[1] = prefix_nr
			// Gateway
			addresses[2] = ipAddrToDecimal(defgateway)

			addressArray := make([][]uint32, 1)
			addressArray[0] = addresses

			connectionSettings[ip4Section][ip4SectionAddresses] = addressArray
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			// Update the connection settings
			err = currentConnections[i].Update(connectionSettings)
			if err != nil {
				log.Print("failed to update connection")
				os.Exit(1)
			}

			// Save the connection settings
			err = currentConnections[i].Save()
			if err != nil {
				log.Print("failed to save setting")
				os.Exit(1)
			}

			fmt.Printf("New IPv4 address set successfully to: %s\n", ip_addr)
			fmt.Printf("New gateway set successfully to: %s\n\n", defgateway)
			return
		}
	}

	log.Print("connection not found in setIPAddr")
	return
}

func setIPMode() {
	log.Print("setting ip mode")

	// Create a new instance of gonetworkmanager.Settings
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		fmt.Print("could net get new settings")
		os.Exit(1)
	}

	// Get the list of all connections
	currentConnections, err := settings.ListConnections()
	if err != nil {
		fmt.Print("could not get settings connections list")
		os.Exit(1)
	}

	// Loop through the connections and find the one with the specified ID
	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			fmt.Print("could not get settings of connection")
			os.Exit(1)
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			connectionSettings[ip4Section][ip4SectionMethod] = ip_mode
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			// order defined by network manager
			addresses := make([]uint32, 3)
			// IP addr
			addresses[0] = ipAddrToDecimal(ip_addr)
			addresses[1] = prefix_nr
			// Gateway
			addresses[2] = ipAddrToDecimal(defgateway)

			addressArray := make([][]uint32, 1)
			addressArray[0] = addresses

			connectionSettings[ip4Section][ip4SectionAddresses] = addressArray
			connectionSettings[ip4Section][ip4SectionNeverDefault] = true
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			// Update the connection settings
			err = currentConnections[i].Update(connectionSettings)
			if err != nil {
				log.Print("failed to update connection")
				os.Exit(1)
			}

			// Save the connection settings
			err = currentConnections[i].Save()
			if err != nil {
				log.Print("failed to save setting")
				os.Exit(1)
			}

			err = settings.ReloadConnections()
			if err != nil {
				log.Print("failed to reload settings")
				os.Exit(1)
			}

			fmt.Printf("IPv4 mode set successfully to %s\n", ip_mode)
			return
		}
	}

	log.Print("connection not found in setIPMode")
	return
}

func refresh_nmcli() {
	log.Print("\nRefreshing nmcli connection on the interface..")
	cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nmcli connection up \""+connectionID+"\"")

	output, err := cmd.CombinedOutput()
	if cmd == nil {
		fmt.Println("failed to execute command")
	}

	cmd = exec.Command("ifconfig", yourInterfaceName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running nmcli:", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func main() {

	get_original_interface_setting()
	setIPAddr()
	setIPMode()
	refresh_nmcli()
}
