package main

import (
	"fmt"
	"net"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to get network interfaces:", err)
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("Interface name: %v\n", iface.Name)
		fmt.Printf("Interface hardware address: %v\n", iface.HardwareAddr)
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Failed to get addresses for interface", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			fmt.Printf("Interface address: %v\n", addr)
		}

		fmt.Printf("\n")
	}
}
