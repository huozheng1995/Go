package main

import (
	"fmt"
	"github.com/edward/mocker/winipcfg"
	"golang.zx2c4.com/wireguard/tun"
	"net"
	"net/netip"
	"strings"
	"time"
)

func CreateInterface(ip string) error {
	fmt.Println("Waiting for network interface to be created... IP: " + ip)

	device, err := tun.CreateTUN("mocker0", 0)
	if err != nil {
		return err
	}
	// Get the LUID to set IP address
	nativeTunDevice := device.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())
	addr, _ := netip.ParsePrefix(ip + "/32")
	err = link.SetIPAddresses([]netip.Prefix{addr})
	if err != nil {
		return err
	}

	//Must stop all operations and wait for the interface to be created!
	time.Sleep(6 * time.Second)
	fmt.Println("Network interface is created! name: mocker0")

	return nil
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
