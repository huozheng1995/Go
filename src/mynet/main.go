package main

import (
	"myutil"
)

var Logger *myutil.MyLogger

func main() {
	Logger = myutil.NewMyLogger("mynet.log")

	config := ParseConfig()

	myNet := NewMyNet(config)
	myNet.Start()
}
