# ipaddresing Go Module

The ipaddressing is a Golang package and it is responsible for handling IP addressing functionality.
It sets the IP address, gateway, and other network settings for a specific interface by using the **gonetworkmanager** package to interact with the network manager (nmcli). 

**func SetupEthInterface(eth_interface_name string, connection_id string)** 
    - This function sets the network interface with the specified name. 
    - The connectionID, which is bound to the interface described via nmcli.
**important** note is that the name of the network interface has to be checked with the "nmcli connection show" command. That way u can link the interface name to the nmcli connection name. (e.g. 'Eth0' is bound to 'Wired connection 1' by the nmcli -> So you pass SetupEthInterface("eth0", "Wired connection 1"))

**func IpAddrToDecimal(ipAddr string) uint32** Takes an IP address as a string and converts it to its decimal representation. It returns the decimal representation as a uint32 value.

**func Get_interface_settings() (error, string, string, string)**
Retrieves the EthInterface name, IP address and gateway of the configured network interface.    

**func SetIPAddr(ip_addr string, prefix_nr uint32, defgateway_addr string) error**
Sets the IP address, prefixnr and gateway address for the setup network interface.

**func SetIPMode(ip_mode string) error** 
Sets the IP mode of the setup interface. It can be either "manual" or "auto".

**func SetDefaultGateway(defgateway_addr string) error**
Sets the default gateway for the specified network interface. 

**func Refresh_nmcli() (error, string)** Refreshes the network connection by executing the nmcli command-line tool. It brings up/updates the **connectionID** setup via the "SetupEthInterface" using the nmcli connection up command. It also retrieves the output of the ifconfig command for the specified network interface for an overview of the network settings.

# Example
inside [/cmd/network/main_network.go](cmd/network/main.go) you can find an example of how to use the ipaddressing package.

# Requirements

go version >1.16

The following packages are required to use the ipaddressing package:
- gonetworkmanager
- dbus

# usage
**GO**
```
go mod init github.com/OpenPixelSystems/go-networkmanager-lib
go run .
```

**Makefile**
```
make all
make tests
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
