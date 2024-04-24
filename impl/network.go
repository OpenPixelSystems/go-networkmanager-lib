package impl

import (
	"net"
	"os"

	"github.com/Wifx/gonetworkmanager/v2"
	"github.com/openpixelsystems/go-networkmanager-lib/interfaces"
)

type Interface struct {
	Interface *net.Interface
}

type NetworkManager struct {
}

func (intf Interface) Addrs() ([]net.Addr, error) {
	return intf.Interface.Addrs()
}

func (NetworkManager) Hostname() (string, error) {
	return os.Hostname()
}

func (NetworkManager) Manager() (gonetworkmanager.NetworkManager, error) {
	return gonetworkmanager.NewNetworkManager()
}

func (NetworkManager) Settings() (gonetworkmanager.Settings, error) {
	return gonetworkmanager.NewSettings()
}

func (NetworkManager) Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

func (NetworkManager) InterfaceByName(networkAdapter string) (interfaces.Interface, error) {
	intf, err := net.InterfaceByName(networkAdapter)
	return &Interface{intf}, err
}
