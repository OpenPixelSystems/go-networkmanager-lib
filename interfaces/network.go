package interfaces

import (
	"net"

	"github.com/Wifx/gonetworkmanager/v2"
)

type Interface interface {
	Addrs() ([]net.Addr, error)
}

type NetworkManager interface {
	Hostname() (string, error)
	Manager() (gonetworkmanager.NetworkManager, error)
	Settings() (gonetworkmanager.Settings, error)
	Interfaces() ([]net.Interface, error)
	InterfaceByName(networkAdapter string) (Interface, error)
}

type Network interface {
	InitializeNetworkState(networkAdapter string)
	RetrieveInternalMode() string
	RetrieveHostname() (string, error)
	RetrieveMACAddress(networkAdapter string) (string, error)
	RetrieveDeviceState(networkAdapter string) gonetworkmanager.NmDeviceState
	RetrieveMode(networkAdapter string) (string, error)
	SetMode(networkAdapter string, mode string) error
	RetrieveIPAddress(networkAdapter string) (string, error)
	SetIPAddress(networkAdapter string, ip string) error
	RetrieveDefaultGateway(networkAdapter string) (string, error)
	SetDefaultGateway(networkAdapter string, gateway string) error
	RetrieveSubnetMask(networkAdapter string) (string, error)
	SetSubnetMask(networkAdapter string, subnet string) error
	RetrieveDNS1(networkAdapter string) (string, error)
	SetDNS1(networkAdapter string, dns string) error
	RetrieveDNS2(networkAdapter string) (string, error)
	SetDNS2(networkAdapter string, dns string) error
	SetIPs(networkAdapter string, mode string, ip string, subnet string, gateway string, dns1 string, dns2 string) error
}
