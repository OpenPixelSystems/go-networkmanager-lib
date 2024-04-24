package bsp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/Wifx/gonetworkmanager/v2"
	"github.com/openpixelsystems/go-networkmanager-lib/impl"
	"github.com/openpixelsystems/go-networkmanager-lib/interfaces"
	"github.com/openpixelsystems/go-networkmanager-lib/utilities"
)

const (
	connectionSection      = "connection"
	connectionSectionID    = "id"
	ip4Section             = "ipv4"
	ip4SectionAddresses    = "addresses"
	ip4SectionAddress      = "address"
	ip4SectionPrefix       = "prefix"
	ip4SectionMethod       = "method"
	ip4SectionGateway      = "gateway"
	ip4SectionNeverDefault = "never-default"
	ip4NameServer          = "dns"
	ip6Section             = "ipv6"
	ip6SectionMethod       = "method"
	connectionID           = "Wired connection 1"

	DNS1Offset = 0
	DNS2Offset = 1
)

type NetworkState struct {
	Mode           string
	IPAddress      string
	SubnetMask     string
	DefaultGateway string
	DNS1           string
	DNS2           string
}

type Network struct {
	NetworkManager interfaces.NetworkManager
	Exec           interfaces.Exec
	Settings       gonetworkmanager.Settings
	Device         gonetworkmanager.Device
	State          NetworkState
}

var DefaultNetwork = Network{
	NetworkManager: &impl.NetworkManager{},
	Exec:           &impl.Exec{},
}

func NewNetwork() *Network {
	network := &Network{}

	network.NetworkManager = DefaultNetwork.NetworkManager
	network.Exec = DefaultNetwork.Exec

	return network
}

func (network *Network) InitializeNetworkState(networkAdapter string) {
	network.State.Mode, _ = network.RetrieveMode(networkAdapter)
	network.State.IPAddress, _ = network.RetrieveIPAddress(networkAdapter)
	network.State.SubnetMask, _ = network.RetrieveSubnetMask(networkAdapter)
	network.State.DefaultGateway, _ = network.RetrieveDefaultGateway(networkAdapter)
	network.State.DNS1, _ = network.RetrieveDNS1(networkAdapter)
	network.State.DNS2, _ = network.RetrieveDNS2(networkAdapter)
}

func (state *NetworkState) stateAllStaticFieldsSet() bool {
	log.Print("State: ", state)
	return state.Mode == "manual" && state.IPAddress != "" && state.SubnetMask != "" && state.DefaultGateway != "" && (state.DNS1 != "" || state.DNS2 != "")
}

func (network *Network) ipAddressToDecimal(ipAddress string) uint32 {
	if ipAddress == "" {
		return 0
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}

	decimal := binary.BigEndian.Uint32(ip.To4())

	// Reverse the endianness
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, decimal)
	reversed := binary.BigEndian.Uint32(bytes)

	return reversed
}

func (network *Network) subnetMaskToCIDR(subnetMask string) uint32 {
	if subnetMask == "" {
		return 0
	}

	ip := net.ParseIP(subnetMask)
	if ip == nil {
		return 0
	}

	cidr, _ := net.IPMask(ip.To4()).Size()

	return uint32(cidr)
}

func (network *Network) setNMGetDevice(networkAdapter string) error {
	if network.Device != nil {
		return nil
	}

	nm, err := network.NetworkManager.Manager()
	if err != nil {
		log.Print(err)
		return err
	}

	device, err := nm.GetDeviceByIpIface(networkAdapter)
	if err != nil {
		log.Print(err)
		return err
	}

	network.Device = device
	return nil
}

func (network *Network) RetrieveInternalMode() string {
	return network.State.Mode
}

func (network *Network) RetrieveHostname() (string, error) {
	return network.NetworkManager.Hostname()
}

func (network *Network) retrieveMACAddress(interfaceStr string) (string, error) {
	log.Print("Retrieving MAC for interface ", interfaceStr)

	ifas, err := network.NetworkManager.Interfaces()
	if err != nil {
		return "", err
	}

	for _, ifa := range ifas {
		if ifa.Name == interfaceStr {
			return ifa.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("interface not found")
}

func (network *Network) RetrieveMACAddress(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	macAddr, err := network.retrieveMACAddress(networkAdapter)
	if err != nil {
		log.Print(err)
		return "", err
	}

	log.Print("MAC address: ", macAddr)
	return macAddr, nil
}

func (network *Network) RetrieveDeviceState(networkAdapter string) gonetworkmanager.NmDeviceState {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return gonetworkmanager.NmDeviceStateUnknown
	}

	state, _ := network.Device.GetPropertyState()

	log.Print("Device state: ", state.String())
	return state
}

func (network *Network) RetrieveMode(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	state, err := network.Device.GetPropertyState()
	if err != nil {
		return "", err
	}

	if state == gonetworkmanager.NmDeviceStateUnavailable {
		return "off", nil
	}

	dhcp4, err := network.Device.GetPropertyDHCP4Config()
	if err != nil {
		return "", err
	}

	mode := "fixed"

	if dhcp4 != nil {
		mode = "dhcp"
	}

	log.Print("IP mode: ", mode)
	return mode, nil
}

func (network *Network) SetMode(networkAdapter string, mode string) error {
	log.Print("setting IP mode: ", mode)
	if mode == "off" {
		network.State.Mode = "off"
	} else if mode == "dhcp" {
		network.State.Mode = "auto"
	} else if mode == "fixed" {
		network.State.Mode = "manual"
		network.State.IPAddress, _ = network.RetrieveIPAddress(networkAdapter)
		network.State.DefaultGateway, _ = network.RetrieveDefaultGateway(networkAdapter)
		network.State.SubnetMask, _ = network.RetrieveSubnetMask(networkAdapter)
		network.State.DNS1, _ = network.RetrieveDNS1(networkAdapter)
		network.State.DNS2, _ = network.RetrieveDNS2(networkAdapter)
	} else {
		return fmt.Errorf("invalid IP mode requested")
	}

	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) CheckInterfacePrivateRange(networkAdapter string) (bool, error) {
	intf, err := network.NetworkManager.InterfaceByName(networkAdapter)
	if err != nil {
		log.Print(err)
		return false, err
	}

	addrs, err := intf.Addrs()
	if err != nil {
		log.Print(err)
		return false, err
	}

	for _, addr := range addrs {
		if addr.(*net.IPNet).IP.IsPrivate() {
			log.Printf("IP addr of %s is private\n", networkAdapter)
			return true, nil
		}
	}

	return false, nil
}

func (network *Network) RetrieveIPAddress(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	ipv4, err := network.Device.GetPropertyIP4Config()
	if err != nil {
		return "", err
	}

	addr, err := ipv4.GetPropertyAddressData()
	if err != nil {
		return "", err
	}

	if len(addr) == 0 {
		return "", fmt.Errorf("no IP address found")
	}

	ip := addr[0].Address

	log.Print("IP address: ", ip)
	return ip, nil
}

func (network *Network) SetIPAddress(networkAdapter string, ip string) error {
	log.Print("setting IP address: ", ip)

	err := utilities.ValidateIPAddress(ip)
	if err != nil {
		return err
	}

	network.State.IPAddress = ip
	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) RetrieveSubnetMask(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	ipv4, err := network.Device.GetPropertyIP4Config()
	if err != nil {
		return "", err
	}

	addr, err := ipv4.GetPropertyAddressData()
	if err != nil {
		return "", err
	}

	if len(addr) == 0 {
		return "", fmt.Errorf("no subnet mask found")
	}

	mask := net.CIDRMask(int(addr[0].Prefix), 32)

	log.Print("Subnet mask: ", net.IP(mask).String())
	return net.IP(mask).String(), nil
}

func (network *Network) SetSubnetMask(networkAdapter string, subnet string) error {
	log.Print("setting subnet mask: ", subnet)

	err := utilities.ValidateIPAddress(subnet)
	if err != nil {
		return err
	}

	network.State.SubnetMask = subnet
	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) RetrieveDefaultGateway(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	ipv4, err := network.Device.GetPropertyIP4Config()
	if err != nil {
		return "", err
	}

	gw, _ := ipv4.GetPropertyGateway()

	log.Print("Default gateway: ", gw)
	return ipv4.GetPropertyGateway()
}

func (network *Network) SetDefaultGateway(networkAdapter string, gateway string) error {
	log.Print("setting default gateway: ", gateway)

	err := utilities.ValidateIPAddress(gateway)
	if err != nil {
		return err
	}

	network.State.DefaultGateway = gateway
	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) RetrieveDNS1(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	ipv4, err := network.Device.GetPropertyIP4Config()
	if err != nil {
		return "", err
	}

	dns, err := ipv4.GetPropertyNameserverData()
	if err != nil {
		return "", err
	}

	nrDNSs := len(dns) / 2

	log.Print("DNS1: ", dns)
	if nrDNSs >= 1 {
		log.Print("DNS1: ", dns[DNS1Offset+nrDNSs].Address)
		return dns[DNS1Offset+nrDNSs].Address, nil
	}

	return "", nil
}

func (network *Network) SetDNS1(networkAdapter string, dns1 string) error {
	log.Print("setting DNS1: ", dns1)

	err := utilities.ValidateIPAddress(dns1)
	if err != nil {
		return err
	}

	network.State.DNS1 = dns1
	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) RetrieveDNS2(networkAdapter string) (string, error) {
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return "", err
	}

	ipv4, err := network.Device.GetPropertyIP4Config()
	if err != nil {
		return "", err
	}

	dns, err := ipv4.GetPropertyNameserverData()
	if err != nil {
		return "", err
	}

	nrDNSs := len(dns) / 2

	log.Print("DNS2: ", dns)
	if nrDNSs >= 2 {
		log.Print("DNS2: ", dns[DNS2Offset+nrDNSs].Address)
		return dns[DNS2Offset+nrDNSs].Address, nil
	}

	return "", nil
}

func (network *Network) SetDNS2(networkAdapter string, dns2 string) error {
	log.Print("setting DNS2: ", dns2)

	err := utilities.ValidateIPAddress(dns2)
	if err != nil {
		return err
	}

	network.State.DNS2 = dns2
	return network.PropagateNetworkSettings(networkAdapter)
}

func (network *Network) SetIPs(networkAdapter string, mode string, ip string, subnet string, gateway string, dns1 string, dns2 string) error {
	log.Print("setting IPs")

	err := network.SetMode(networkAdapter, mode)
	if err != nil {
		return err
	}

	if mode == "fixed" {
		err := utilities.ValidateIPAddress(ip)
		if err != nil {
			return err
		}

		err = utilities.ValidateIPAddress(subnet)
		if err != nil {
			return err
		}

		err = utilities.ValidateIPAddress(gateway)
		if err != nil {
			return err
		}

		err = utilities.ValidateIPAddress(dns1)
		if err != nil {
			return err
		}

		if dns2 != "" { // DNS2 is optional in the group, and can be left empty
			err = utilities.ValidateIPAddress(dns2)
			if err != nil {
				return err
			}
		}

		network.State.IPAddress = ip
		network.State.SubnetMask = subnet
		network.State.DefaultGateway = gateway
		network.State.DNS1 = dns1
		network.State.DNS2 = dns2

		return network.PropagateNetworkSettings(networkAdapter)
	}

	return nil
}

func (network *Network) PropagateNetworkSettings(networkAdapter string) error {
	log.Print("Propageting network settings")
	err := network.setNMGetDevice(networkAdapter)
	if err != nil {
		return err
	}

	settings, err := network.NetworkManager.Settings()
	if err != nil {
		log.Print("could not get settings")
		return err
	}

	network.Settings = settings // For testing only!
	currentConnections, err := settings.ListConnections()
	if err != nil {
		log.Print("could not get settings")
		return err
	}

	for i := range currentConnections {
		connectionSettings, err := currentConnections[i].GetSettings()
		if err != nil {
			log.Print("could not get settings of connection")
			return err
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			log.Print(connectionSettings)

			// Set IP mode
			connectionSettings[ip4Section][ip4SectionMethod] = network.State.Mode

			// Ignore IPv6
			connectionSettings[ip6Section] = make(map[string]interface{})
			connectionSettings[ip6Section][ip6SectionMethod] = "ignore"

			log.Print(network.State.Mode)

			// Network State Handling
			currentMode, _ := network.RetrieveMode(networkAdapter)
			if currentMode == "off" {
				log.Print("mode is off, not propagating settings")
				return nil
			} else if network.State.Mode == "auto" {
				// Clear addresses
				empty := make([][]uint32, 1)
				connectionSettings[ip4Section][ip4SectionAddresses] = empty

				var dnss []uint32
				connectionSettings[ip4Section][ip4NameServer] = dnss

				// Clear all fields, except the mode one
				network.State = NetworkState{
					Mode: "auto",
				}
			} else if network.State.Mode == "manual" {
				// If not all fields set, return
				log.Print(network.State)
				if !network.State.stateAllStaticFieldsSet() {
					log.Print("Not all static parameters set...")
					return nil
				}

				addressData := make([]map[string]interface{}, 1)
				addressData[0] = make(map[string]interface{})

				subnetMaskCIDR := network.subnetMaskToCIDR(network.State.SubnetMask)

				addressData[0][ip4SectionPrefix] = subnetMaskCIDR
				addressData[0][ip4SectionAddress] = network.State.IPAddress

				// order defined by network manager
				addresses := make([]uint32, 3)

				// IP address
				addresses[0] = network.ipAddressToDecimal(network.State.IPAddress) // e.g., 192.168.12.2

				// Subnet mask (CIDR notation)
				addresses[1] = subnetMaskCIDR // e.g., 24

				// Default gateway
				addresses[2] = network.ipAddressToDecimal(network.State.DefaultGateway) // e.g., 192.168.12.1

				// Set all addresses
				addressArray := make([][]uint32, 1)
				addressArray[0] = addresses
				connectionSettings[ip4Section][ip4SectionAddresses] = addressArray

				// Set DNSs

				var dnss []uint32
				dns1 := network.ipAddressToDecimal(network.State.DNS1)
				if dns1 != 0 {
					dnss = append(dnss, dns1)
				}

				dns2 := network.ipAddressToDecimal(network.State.DNS2)
				if dns2 != 0 {
					dnss = append(dnss, dns2)
				}

				log.Print("dnss: ", dnss)
				if len(dnss) != 0 {
					connectionSettings[ip4Section][ip4NameServer] = dnss
				}
			} else {
				log.Print(network.State.Mode, "is invalid")
				return fmt.Errorf("invalid mode")
			}

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

			err = settings.ReloadConnections()
			if err != nil {
				log.Print("failed to reload settings")
				return err
			}

			cmd := network.Exec.Command("/bin/sh", "-c", "/usr/bin/nmcli connection up \"Wired connection 1\"")
			if cmd == nil {
				return fmt.Errorf("failed to execute command")
			}

			out, err := cmd.Output()
			if err != nil {
				return err
			}

			log.Print(string(out))
			return nil
		}
	}

	log.Print("connection not found!")
	return fmt.Errorf("connection not found")
}
