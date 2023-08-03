package logger

import (
	stlog "log"
	"os"
	"strconv"
)

var log *stlog.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func InitLog(destination string) {
	log = stlog.New(fileLog(destination), "", stlog.LstdFlags)
}

func Log(code string, message string) {
	log.Printf("[%v] %v\n", code, message)
}

func LogBytes(code string, message string, arr []byte, printDetails bool) {
	if printDetails {
		log.Printf("[%v] %v\n%v\n", code, message, arr)
	} else {
		log.Printf("[%v] %v\n", code, message)
	}
}

func LogWarn(code string, message string) {
	log.Printf("[%v] [WARN]%v\n", code, message)
}

func LogError(code string, message string) {
	log.Printf("[%v] [ERROR]%v\n", code, message)
}

func GetCode(connId int, socket string) string {
	return "Conn" + strconv.Itoa(connId) + "-" + socket
}
