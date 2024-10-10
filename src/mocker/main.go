package main

import (
	"encoding/json"
	"myutil"
	"os"
	"os/exec"
	"strconv"
)

var Logger *myutil.MyLogger

func main() {
	Logger = myutil.NewMyLogger("mocker.log")

	config := ParseConfig()
	config.MockerName = "mocker0"

	go createNetwork(config.MockerIP, config.MockerName)

	mocker := NewMocker(config)
	mocker.Start()
}

func createNetwork(ServerIP string, networkName string) {
	configObj := struct {
		NetworksToAdd []myutil.Network `json:"NetworksToAdd"`
	}{
		NetworksToAdd: []myutil.Network{
			{Name: networkName, IPv4Address: ServerIP, SubnetMask: 32},
		},
	}

	configStr, _ := json.Marshal(configObj)

	file, _ := os.Create("mynet/config.json")
	_, err := file.Write(configStr)
	if err != nil {
		Logger.LogError("MyNet", "Error writing config file, error: "+err.Error())
		panic(err)
	}
	file.Close()

	cmd := exec.Command("mynet\\start.bat")
	err = cmd.Start()
	if err != nil {
		Logger.LogError("MyNet", "Error running mynet.exe, error: "+err.Error())
	}

	Logger.Log("MyNet", "mynet.exe is started, process id: "+strconv.Itoa(cmd.Process.Pid))
}
