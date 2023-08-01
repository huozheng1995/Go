package main

import (
	"errors"
	"github.com/edward/mocker/winipcfg"
	"golang.zx2c4.com/wireguard/tun"
	"net"
	"net/netip"
	"strings"
)

func CreateNetworkInterface(ip string) (*tun.NativeTun, error) {
	device, err := tun.CreateTUN("mocker", 0)
	if err != nil {
		return nil, err
	}
	// Get the LUID to set IP address
	nativeTunDevice := device.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())
	addr, err := netip.ParsePrefix(ip + "/24")
	if err != nil {
		return nil, err
	}
	//err = link.SetIPAddresses([]netip.Prefix{addr})

	//addr1, _ := netip.ParsePrefix("172.16.88.1/32")
	err = link.SetIPAddresses([]netip.Prefix{addr})
	//link.AddRoute(addr, addr.Addr(), 0)

	if err != nil {
		return nil, err
	}

	return nativeTunDevice, nil
}

func GetNetworkInterface(ip string) (*net.Interface, error) {
	nilResult := net.Interface{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return &nilResult, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			LogError("Failed to get addresses for interface:"+iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), ip) {
				return &iface, nil
			}
		}
	}

	return &nilResult, errors.New("There is no matching network interface, ip: " + ip)
}
