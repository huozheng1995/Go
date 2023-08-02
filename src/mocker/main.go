package main

import "github.com/edward/mocker/logger"

type Settings struct {
	MockerPort                   int    `json:"MockerPort"`
	ServerIP                     string `json:"ServerIP"`
	ServerPort                   int    `json:"ServerPort"`
	PrintDetails                 bool   `json:"PrintDetails"`
	VirtualNetworkInterfaceMode  bool   `json:"VirtualNetworkInterfaceMode"`
	LocalNetworkInterfaceAddress string `json:"LocalNetworkInterfaceAddress"`
}

func main() {
	logger.InitLog("mocker.log")

	settings := Settings{
		MockerPort:                   3306,
		ServerIP:                     "172.16.85.139",
		ServerPort:                   3306,
		PrintDetails:                 true,
		VirtualNetworkInterfaceMode:  true,
		LocalNetworkInterfaceAddress: "172.16.85.1",
	}
	m := NewMocker(settings)
	//AddReqDataResData(m, []byte("abc\n"), []byte("def\n"))
	//AddReqDataResData(m, HexFileToBytes("D:\\log\\mock\\columns_req.txt"), HexFileToBytes("D:\\log\\mock\\columns_res.txt"))
	//AddReqDataResDataFiles(m, HexFileToBytes("C:\\Users\\User\\Downloads\\mocker\\req.txt"),
	//	"C:\\Users\\User\\Downloads\\mocker\\1.log",
	//	"C:\\Users\\User\\Downloads\\mocker\\2.log",
	//	"C:\\Users\\User\\Downloads\\mocker\\3.log",
	//	"C:\\Users\\User\\Downloads\\mocker\\4.log")
	//AddResLenResData(m, 1945, HexFileToBytes("C:\\Users\\User\\Downloads\\SYNC-2347\\D_R_S\\log.log"))

	m.Start()
}
