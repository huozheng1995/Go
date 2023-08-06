package main

import (
	"errors"
	"fmt"
	"github.com/edward/mocker/winipcfg"
	"golang.zx2c4.com/wireguard/tun"
	"net"
	"net/netip"
	"strings"
	"time"
)

func CreateInterface(ip string) (*tun.NativeTun, error) {
	addr, err := netip.ParsePrefix(ip + "/32")
	if err != nil {
		return nil, err
	}

	device, err := tun.CreateTUN("mocker0", 0)
	if err != nil {
		return nil, err
	}
	// Get the LUID to set IP address
	nativeTunDevice := device.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())
	err = link.SetIPAddresses([]netip.Prefix{addr})
	if err != nil {
		return nil, err
	}

	//Must stop all operations and wait for the interface to be created!
	fmt.Println("Waiting for network interface to be created...")
	//Logger.Log("Main", "Waiting for network interface to be created...")
	time.Sleep(3 * time.Second)
	msg2 := "Network interface is created! name: mocker0, address: " + addr.String()
	fmt.Println(msg2)
	Logger.Log("Main", msg2)

	return nativeTunDevice, nil
}

func GetInterface(ip string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			Logger.LogError("Main", "Failed to get addresses for interface:"+iface.Name+", error: "+err.Error())
			continue
		}
		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), ip) {
				return &iface, nil
			}
		}
	}

	return nil, errors.New("There is no matching network interface, ip: " + ip)
}

func GetInterfaceAddr(interfaceName string) (*net.IP, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		Logger.LogError("Main", "Failed to get addresses for interface:"+iface.Name+", error: "+err.Error())
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return &ipnet.IP, nil
			}
		}
	}

	return nil, errors.New("There is no Addr for interface, name: " + interfaceName)
}
