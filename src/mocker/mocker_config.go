package main

import (
	"encoding/json"
	"os"
)

type MockerConfig struct {
	ServerIP         string `json:"ServerIP"`
	ServerPort       int    `json:"ServerPort"`
	MockerName       string `json:"MockerName"`
	MockerIP         string `json:"MockerIP"`
	MockerPort       int    `json:"MockerPort"`
	DefNetInterface  string `json:"DefNetInterface"`
	PrintDetails     bool   `json:"PrintDetails"`
	MockDataLocation string `json:"MockDataLocation"`
	MockDataGroup1   []struct {
		RequestFile   string   `json:"RequestFile"`
		ResponseFiles []string `json:"ResponseFiles"`
	} `json:"MockDataGroup1"`
	MockDataGroup2 []struct {
		ResponseDataLength int      `json:"ResponseDataLength"`
		ResponseFiles      []string `json:"ResponseFiles"`
	} `json:"MockDataGroup2"`
	MockDataGroup3 []struct {
		RequestFile string `json:"RequestFile"`
		Delay       int    `json:"Delay"`
	} `json:"MockDataGroup3"`
}

func ParseConfig() *MockerConfig {
	file, err := os.Open("config.json")
	if err != nil {
		Logger.LogError("Main", "Error opening config file:"+err.Error())
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := MockerConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		Logger.LogError("Main", "Error decoding config file:"+err.Error())
		panic(err)
	}

	return &config
}
