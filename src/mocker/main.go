package main

type Settings struct {
	MockerPort                   int    `json:"MockerPort"`
	ServerIP                     string `json:"ServerIP"`
	ServerPort                   int    `json:"ServerPort"`
	PrintDetails                 bool   `json:"PrintDetails"`
	VirtualNetworkInterfaceMode  bool   `json:"VirtualNetworkInterfaceMode"`
	LocalNetworkInterfaceAddress string `json:"LocalNetworkInterfaceAddress"`
}

func main() {
	InitLog("mocker.log")

	settings := Settings{
		MockerPort:                   3306,
		ServerIP:                     "172.16.85.140",
		ServerPort:                   3306,
		PrintDetails:                 true,
		VirtualNetworkInterfaceMode:  true,
		LocalNetworkInterfaceAddress: "172.16.85.1",
	}
	m := NewMocker(settings)
	AddReqDataResData(m, []byte("a"), []byte("b"))
	AddReqDataResData(m, []byte("c"), []byte("d"))
	AddReqDataResData(m, []byte("abc\n"), []byte("def\n"))
	//AddReqDataResData(m, HexFileToBytes("C:\\Users\\33907\\Downloads\\req.txt"), HexFileToBytes("C:\\Users\\33907\\Downloads\\res.txt"))
	//AddResLenResData(m, 1945, HexFileToBytes("C:\\Users\\User\\Downloads\\SYNC-2347\\D_R_S\\log.log"))

	m.Start()
}
