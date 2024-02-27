package main

import (
	"log"

	"openpixelsystems.org/go-networkmanager-lib/ethInterfacing"
)

func Setup_EthInterface() (error, string) {
	err := ethInterfacing.Get_original_interface_setting()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = ethInterfacing.SetIPAddr()
	if err != nil {
		log.Fatal(err)
		return err, ""
	}

	err = ethInterfacing.SetIPMode()
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
	}
	log.Print(strout)
}
