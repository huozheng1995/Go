package main

import (
	"errors"
	"github.com/edward/mocker/logger"
	"github.com/edward/mocker/winipcfg"
	"golang.zx2c4.com/wireguard/tun"
	"net"
	"net/netip"
	"strings"
)

func CreateNetworkInterface(ip string) (*tun.NativeTun, error) {
	device, err := tun.CreateTUN("mocker0", 0)
	if err != nil {
		return nil, err
	}
	// Get the LUID to set IP address
	nativeTunDevice := device.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())

	addr, err := netip.ParsePrefix(ip + "/32")
	if err != nil {
		return nil, err
	}

	err = link.SetIPAddresses([]netip.Prefix{addr})
	if err != nil {
		return nil, err
	}
	logger.Log("Network Interface is created! name: mocker0, address: " + addr.String())

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
			logger.LogError("Failed to get addresses for interface:"+iface.Name, err)
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
