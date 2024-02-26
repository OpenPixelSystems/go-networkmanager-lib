# ipaddresing Go Module

The ipaddressing is a Golang package and it is responsible for handling IP addressing functionality.
It sets the IP address, gateway, and other network settings for a specific interface by using the **gonetworkmanager** package to interact with the network manager (nmcli). 

ipAddrToDecimal(ipAddr string) uint32: Takes an IP address as a string and converts it to its decimal representation. It returns the decimal representation as a uint32 value.

get_original_interface_setting(): Retrieves the original settings of the network interface.

setIPAddr(): Sets the IP address and gateway for the specified network interface. It uses the values defined in the constants ip_addr and defgateway to set the IP address and gateway respectively. 

setIPMode(): Sets the IP mode for the specified network interface. It uses the value defined in the constant ip_mode to set the IP mode. 

refresh_nmcli(): Refreshes the network connection using the nmcli command-line tool. It brings up the "Wired connection 1" connection using the nmcli connection up command. It also retrieves the output of the ifconfig command for the specified network interface. 
**important** note is that the name of the network interface has to be checked with the "nmcli connection show" command. That way u can link the interface name to the nmcli connection name.

# Requirements

go version >1.16

The following packages are required to use the ipaddressing package:
- gonetworkmanager
- dbus

# usage

```
go run .
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
