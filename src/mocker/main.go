package main

type Settings struct {
	MockerPort   int    `json:"MockerPort"`
	ServerIP     string `json:"ServerIP"`
	ServerPort   int    `json:"ServerPort"`
	PrintDetails bool   `json:"PrintDetails"`
}

func main() {
	InitLog("mocker.log")

	settings := Settings{
		MockerPort:   11111,
		ServerIP:     "adb-8439982502599436.16.azuredatabricks.net",
		ServerPort:   443,
		PrintDetails: true,
	}
	m := NewMocker(settings)
	AddReqDataResData(m, []byte("abc\n"), []byte("def\n"))
	//AddReqDataResData(m, HexFileToBytes("C:\\Users\\33907\\Downloads\\req.txt"), HexFileToBytes("C:\\Users\\33907\\Downloads\\res.txt"))
	//AddResLenResData(m, 1945, HexFileToBytes("C:\\Users\\User\\Downloads\\SYNC-2347\\D_R_S\\log.log"))

	m.Start()
}
