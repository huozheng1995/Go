package main

import (
	"myutil"
)

var Logger *myutil.MyLogger

func main() {
	Logger = myutil.NewMyLogger("mocker.log")

	config := ParseConfig()

	m := NewMocker(config)
	m.Start()
}
