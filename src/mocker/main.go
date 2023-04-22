package main

const (
	IP   = "172.16.85.140"
	Port = 3306
)

func main() {
	InitLog("mocker.log")

	m := NewMocker(IP, Port)
	//AddReqDataResData(m, []byte("abc\n"), []byte("def\n"))
	//AddReqDataResData(m, HexFileToBytes("C:\\Users\\33907\\Downloads\\req.txt"), HexFileToBytes("C:\\Users\\33907\\Downloads\\res.txt"))
	//AddResLenResData(m, 356, HexFileToBytes("C:\\Users\\33907\\Downloads\\res.txt"))

	m.Start()
}
