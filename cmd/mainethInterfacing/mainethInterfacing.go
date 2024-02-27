package main

import (
	"log"
	"os"

	"openpixelsystems.org/go-networkmanager-lib/ethInterfacing"
)

func Setup_EthInterface() (error, string) {
	err := ethInterfacing.SetIPAddr()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = ethInterfacing.SetIPMode()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = ethInterfacing.SetDefaultGateway()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err, stroutput := ethInterfacing.Refresh_nmcli()
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

	err, ethface, ethaddr, ethgateway := ethInterfacing.Get_interface_settings()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Interface Settings: %s, IP: %s, Gateway: %s", ethface, ethaddr, ethgateway)

}
