package myutil

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Network struct {
	Name        string `json:"Name"`
	IPv4Address string `json:"IPv4Address"`
	SubnetMask  int    `json:"SubnetMask"`
}

func CreateInterfaceManual(ip string) error {
	for {
		fmt.Println("Waiting for network interface to be created... IP: " + ip)
		iface, err := GetInterfaceByIP(ip)
		if err != nil {
			return err
		}
		if iface != nil {
			fmt.Println("Network interface is created! name: " + iface.Name)
			return nil
		}
		time.Sleep(3 * time.Second)
	}
}

func GetInterfaceByIP(ip string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), ip) {
				return &iface, nil
			}
		}
	}

	return nil, nil
}

func GetInterfaceIPByName(name string) (*net.IP, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return &ipnet.IP, nil
			}
		}
	}

	return nil, nil
}

func GetIPByDomain(domain string) (*net.IP, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return &ip, nil
		}
	}

	return nil, nil
}

func FindInterfaceInRouteTable(ipStr string, interfaceToIgnore string, defNetInterface string) (*net.IPNet, error) {
	ip := net.ParseIP(ipStr)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == interfaceToIgnore {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.Contains(ip) {
					return ipnet, nil
				} else if iface.Name == defNetInterface && ipnet.IP.To4() != nil {
					return ipnet, nil
				}
			}
		}
	}

	return nil, nil
}
