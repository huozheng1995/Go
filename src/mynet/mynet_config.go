package main

import (
	"encoding/json"
	"myutil"
	"os"
)

type MyNetConfig struct {
	NetworksToAdd []myutil.Network `json:"NetworksToAdd"`
}

func ParseConfig() *MyNetConfig {
	file, err := os.Open("config.json")
	if err != nil {
		Logger.LogError("Main", "Error opening config file:"+err.Error())
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := MyNetConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		Logger.LogError("Main", "Error decoding config file:"+err.Error())
		panic(err)
	}

	return &config
}
