package main

import (
	"strconv"
)

const (
	IP   = "172.16.85.140"
	Port = 3306
)

func main() {
	InitLog("mocker.log")

	m := NewMocker(IP, Port)

	AddReqRes(m, []byte("abc\n"), []byte("def\n"))
	AddReqResFromFile(m, "C:\\Users\\33907\\Downloads\\req.txt", "C:\\Users\\33907\\Downloads\\res.txt")

	Log("Mocker started on port " + strconv.Itoa(Port))
	m.Start()
}
