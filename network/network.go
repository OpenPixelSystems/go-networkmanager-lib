package network

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/Wifx/gonetworkmanager/v2"
)

const (
	connectionSection      = "connection"
	connectionSectionID    = "id"
	ip4Section             = "ipv4"
	ip4SectionAddresses    = "addresses"
	ip4SectionAddress      = "address"
	ip4SectionPrefix       = "prefix"
	ip4SectionMethod       = "method"
	ip4SectionNeverDefault = "never-default"
	ip4SectionGateway      = "gateway"
	ip6Section             = "ipv6"
	ip6SectionMethod       = "method"
)

var EthInterfaceName string
var connectionID string

func SetupEthInterface(eth_interface_name string, connection_id string) {
	EthInterfaceName = eth_interface_name
	connectionID = connection_id
}

func IpAddrToDecimal(ipAddr string) uint32 {
	ip := net.ParseIP(ipAddr)
	decimal := binary.BigEndian.Uint32(ip.To4())

	// Reverse the endianness
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, decimal)
	reversed := binary.BigEndian.Uint32(bytes)
	return reversed
}

func Get_interface_settings() (error, string, string, string) {
	log.Print("getting original interface settings")
	// Create a new instance of gonetworkmanager
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		fmt.Println(err.Error())
		return err, "", "", ""
	}

	// Get the list of all network devices
	devices, err := nm.GetDevices()
	if err != nil {
		fmt.Println(err.Error())
		return err, "", "", ""
	}

	// Find the device with the specified interface name
	var selectedDevice gonetworkmanager.Device
	var dev_inf_name string
	for _, device := range devices {
		interfaceName, _ := device.GetPropertyInterface()
		if interfaceName == EthInterfaceName { // Change to the interface name you want to use
			selectedDevice = device
			dev_inf_name = interfaceName
			break
		}
	}

	if selectedDevice == nil {
		fmt.Println("Selected device not found")
		return err, "", "", ""
	}

	// Get the original IPv4 settings of the selected device
	IPv4Settings, err := selectedDevice.GetPropertyIP4Config()
	if err != nil {
		fmt.Println("Failed to get original IPv4 settings:", err)
		return err, "", "", ""
	}

	// Extract IPv4 addresses
	IPv4Addresses, err := IPv4Settings.GetPropertyAddressData()
	if err != nil {
		fmt.Println("Failed to get original IPv4 addresses:", err)
		return err, "", "", ""
	}

	// Extract the gateway
	IPv4Gateway, err := IPv4Settings.GetPropertyGateway()
	if err != nil {
		fmt.Println("Failed to get original gateway:", err)
		return err, "", "", ""
	}

	IPv4Address := IPv4Addresses[0].Address
	IPv4Prefix := IPv4Addresses[0].Prefix

	fmt.Printf("Got interface name: %s\n", dev_inf_name)
	fmt.Printf("Got IPv4 address: %s/%d\n", IPv4Address, IPv4Prefix)
	fmt.Printf("Got gateway: %s\n\n", IPv4Gateway)
	return nil, dev_inf_name, IPv4Address, IPv4Gateway
}

func SetDefaultGateway(defgateway_addr string) error {
	log.Print("setting default gateway")

	// Create a new instance of gonetworkmanager.Settings
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		fmt.Print("could not get new settings")
		return err
	}

	// Get the list of all connections
	currentConnections, err := settings.ListConnections()
	if err != nil {
		fmt.Print("could not get settings connections list")
		return err
	}

	// Loop through the connections and find the one with the specified ID
	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			fmt.Print("could not get settings of connection")
			return err
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			addressData := make([]map[string]interface{}, 1)
			addressData[0] = make(map[string]interface{})
			addressData[0][ip4SectionGateway] = defgateway_addr

			// order defined by network manager
			addresses := make([][]uint32, 1)
			addresses = connectionSettings[ip4Section][ip4SectionAddresses].([][]uint32)

			// Gateway
			addresses[0][2] = IpAddrToDecimal(defgateway_addr)

			addressArray := make([][]uint32, 1)
			addressArray[0] = addresses[0]

			connectionSettings[ip4Section][ip4SectionAddresses] = addressArray
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			err = currentConnections[i].Update(connectionSettings)
			if err != nil {
				log.Print("failed to update connection")
				return err
			}

			err = currentConnections[i].Save()
			if err != nil {
				log.Print("failed to save setting")
				return err
			}

			log.Print("connection reloaded")
			return nil
		}
	}

	return fmt.Errorf("connection not found in setDefaultGateway")
}

func SetIPAddr(ip_addr string, prefix_nr uint32, defgateway_addr string) error {
	log.Print("setting ip address")

	// Create a new instance of gonetworkmanager.Settings
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		fmt.Print("could not get new settings")
		return err
	}

	// Get the list of all connections
	currentConnections, err := settings.ListConnections()
	if err != nil {
		fmt.Print("could not get settings connections list")
		return err
	}

	// Loop through the connections and find the one with the specified ID
	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			fmt.Print("could not get settings of connection")
			return err
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
			addresses[0] = IpAddrToDecimal(ip_addr)
			addresses[1] = prefix_nr
			// Gateway
			addresses[2] = IpAddrToDecimal(defgateway_addr)

			addressArray := make([][]uint32, 1)
			addressArray[0] = addresses

			connectionSettings[ip4Section][ip4SectionAddresses] = addressArray
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			// Update the connection settings
			err = currentConnections[i].Update(connectionSettings)
			if err != nil {
				log.Print("failed to update connection")
				return err
			}

			// Save the connection settings
			err = currentConnections[i].Save()
			if err != nil {
				log.Print("failed to save setting")
				return err
			}

			fmt.Printf("New IPv4 address set successfully to: %s\n", ip_addr)
			fmt.Printf("New gateway set successfully to: %s\n\n", defgateway_addr)
			return nil
		}
	}

	return fmt.Errorf("connection not found in setIPAddr")
}

func SetIPMode(ip_mode string) error {
	log.Print("setting ip mode")

	if ip_mode != "auto" && ip_mode != "manual" {
		return fmt.Errorf("invalid ip mode")
	}
	// Create a new instance of gonetworkmanager.Settings
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		fmt.Print("could net get new settings")
		return err
	}

	// Get the list of all connections
	currentConnections, err := settings.ListConnections()
	if err != nil {
		fmt.Print("could not get settings connections list")
		return err
	}

	// Loop through the connections and find the one with the specified ID
	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			fmt.Print("could not get settings of connection")
			return err
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			connectionSettings[ip4Section][ip4SectionMethod] = ip_mode

			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"
			connectionSettings[ip4Section][ip4SectionNeverDefault] = false
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			// Update the connection settings
			err = currentConnections[i].Update(connectionSettings)
			if err != nil {
				log.Print("failed to update connection")
				return err
			}

			// Save the connection settings
			err = currentConnections[i].Save()
			if err != nil {
				log.Print("failed to save setting")
				return err
			}

			err = settings.ReloadConnections()
			if err != nil {
				log.Print("failed to reload settings")
				return err
			}

			fmt.Printf("IPv4 mode set successfully to %s\n", ip_mode)
			return nil
		}
	}

	return fmt.Errorf("connection not found in setIPMode")
}

func Refresh_nmcli() (error, string) {
	log.Print("\nRefreshing nmcli connection on the interface..")
	cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nmcli connection up \""+connectionID+"\"")

	output, err := cmd.CombinedOutput()
	if cmd == nil {
		fmt.Println("failed to execute command")
		return err, ""
	}

	cmd = exec.Command("ifconfig", EthInterfaceName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running nmcli:", err)
		return err, ""
	}
	fmt.Println(string(output))
	return nil, string(output)
}
