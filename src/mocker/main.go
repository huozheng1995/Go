package main

import (
	"myutil"
)

var Logger *myutil.MyLogger

func main() {
	Logger = myutil.NewMyLogger("mocker.log")
	defer Logger.Close()

	config := ParseConfig()

	m := NewMocker(config)
	m.Start()
}
