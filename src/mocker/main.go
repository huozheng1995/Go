package main

const (
	MockerPort = 3306
	ServerIP   = "172.16.85.138"
	ServerPort = 3306
)

func main() {
	InitLog("mocker.log")

	m := NewMocker(MockerPort, ServerIP, ServerPort)
	AddReqDataResData(m, []byte("abc\n"), []byte("def\n"))
	//AddReqDataResData(m, HexFileToBytes("C:\\Users\\33907\\Downloads\\req.txt"), HexFileToBytes("C:\\Users\\33907\\Downloads\\res.txt"))
	//AddResLenResData(m, 1945, HexFileToBytes("C:\\Users\\User\\Downloads\\SYNC-2347\\D_R_S\\log.log"))
	AddResLenResData(m, 1945, HexFileToBytes("C:\\Users\\User\\Downloads\\SYNC-2347\\D_R_S\\log2.log"))

	m.Start()
}
