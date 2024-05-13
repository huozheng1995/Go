package main

import (
	"fmt"
	"github.com/edward/mynet/winipcfg"
	"golang.org/x/net/ipv4"
	"golang.zx2c4.com/wireguard/tun"
	"myutil"
	"net/netip"
	"strconv"
	"time"
)

type MyNet struct {
	NetworksToAdd *myutil.Set
	MyNetConfig   *MyNetConfig
}

func NewMyNet(config *MyNetConfig) *MyNet {
	networksToAdd := myutil.NewSet()
	for _, network := range config.NetworksToAdd {
		networksToAdd.Add(&myutil.Network{
			Name:        network.Name,
			IPv4Address: network.IPv4Address,
			SubnetMask:  network.SubnetMask,
		})
	}
	return &MyNet{
		NetworksToAdd: networksToAdd,
		MyNetConfig:   config,
	}
}

func (m *MyNet) Start() {
	networksToAdd := m.NetworksToAdd.Elements()
	for _, item := range networksToAdd {
		network := item.(*myutil.Network)
		var device tun.Device
		var err error
		device, err = CreateNetwork(network.Name, network.IPv4Address, network.SubnetMask)
		if err != nil {
			Logger.LogError("Main", "Error creating network interface, error: "+err.Error())
			panic(err)
		}

		go LogNetworkData(device)
	}

	time.Sleep(24 * time.Hour)
}

func CreateNetwork(name string, ip string, subnetMask int) (tun.Device, error) {
	fmt.Println("Creating network interface... IP: " + ip)
	addr, err := netip.ParsePrefix(ip + "/" + strconv.Itoa(subnetMask))
	if err != nil {
		return nil, err
	}

	device, err := tun.CreateTUN(name, 0)
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

	time.Sleep(3 * time.Second)
	fmt.Println("Network interface is created! name: " + name)

	return device, nil
}

func LogNetworkData(device tun.Device) error {
	const ProtocolTCP = 6
	const MaxIPPacketLen = 1500
	buf1 := make([]byte, MaxIPPacketLen)
	buf2 := make([]byte, MaxIPPacketLen)
	buf3 := make([]byte, MaxIPPacketLen)
	buf4 := make([]byte, MaxIPPacketLen)
	bufs := [][]byte{buf1, buf2, buf3, buf4}
	sizes := make([]int, len(bufs))
	offset := 0

	for {
		filledBufNum, err := device.Read(bufs, sizes, offset)
		if err != nil {
			Logger.LogError("Main", "Error reading data from device:"+err.Error())
			panic(err)
		}
		for i := 0; i < filledBufNum; i++ {
			buf := bufs[i]
			size := sizes[i]
			header, err := ipv4.ParseHeader(buf[offset : offset+size])
			if err != nil {
				Logger.LogError("Main", "Error parsing IP packet header")
				Logger.LogBytes("Main", "Data: ", buf[offset:offset+size], true)
				continue
			}
			if header.Protocol == ProtocolTCP {
				Logger.LogBytes("Main", "Read TCP packet, length: "+strconv.Itoa(size-header.Len), buf[offset+header.Len:offset+size], true)
			} else {
				Logger.Log("Main", "Unknown protocol, IP header: "+header.String())
			}
		}
	}
}
