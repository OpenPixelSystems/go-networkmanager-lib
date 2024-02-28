package main

import (
	"log"
	"os"

	"github.com/OpenPixelSystems/go-networkmanager-lib/network"
)

const (
	ip_addr    = "10.0.3.30" // IP address you want to set
	defgateway = "10.0.0.1"  // Gateway you want to set
	prefix_nr  = 20          // Prefix number
	ip_mode    = "manual"    // Mode you want to set

	ethInterfaceName = "eth0"
	ConnectionID     = "Wired connection 1"
)

func Setup_EthInterface() (error, string) {
	network.SetupEthInterface(ethInterfaceName, ConnectionID)

	err := network.SetIPAddr(ip_addr, prefix_nr, defgateway)
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = network.SetIPMode(ip_mode)
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = network.SetDefaultGateway(defgateway)
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err, stroutput := network.Refresh_nmcli()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}
	log.Print(stroutput)
	return err, stroutput
}

func main() {
	err, strout := Setup_EthInterface()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Print(strout)

	err, ethface, ethaddr, ethgateway := network.Get_interface_settings()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Interface Settings: %s, IP: %s, Gateway: %s", ethface, ethaddr, ethgateway)

}
