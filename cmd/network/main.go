package main

import (
	"flag"

	"github.com/openpixelsystems/go-networkmanager-lib/impl"
	bsp "github.com/openpixelsystems/go-networkmanager-lib/network"
)

func main() {
	exec := &impl.Exec{}
	nm := &impl.NetworkManager{}
	network := &bsp.Network{NetworkManager: nm, Exec: exec}

	networkAdapter := flag.String("adapter", "", "The network adapter to use")
	flag.Parse()
	network.InitializeNetworkState(*networkAdapter)
}
