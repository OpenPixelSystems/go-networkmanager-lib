package main

import (
	"log"

	"openpixelsystems.org/go-networkmanager-lib/ethInterfacing"
)

func main() {
	err := ethInterfacing.Get_original_interface_setting()
	if err != nil {
		log.Fatal(err)
	}

	err = ethInterfacing.SetIPAddr()
	if err != nil {
		log.Fatal(err)
	}

	err = ethInterfacing.SetIPMode()
	if err != nil {
		log.Fatal(err)
	}
}
