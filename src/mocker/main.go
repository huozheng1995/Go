package main

import "github.com/edward/mocker/logger"

func main() {
	logger.InitLog("mocker.log")
	config := ParseConfig()

	m := NewMocker(config)
	m.Start()
}
