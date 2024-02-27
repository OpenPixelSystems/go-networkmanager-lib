package main

import (
	"log"
	"os"

	"github.com/OpenPixelSystems/go-networkmanager-lib/network"
)

func Setup_EthInterface() (error, string) {
	err := network.SetIPAddr()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = network.SetIPMode()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = network.SetDefaultGateway()
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
