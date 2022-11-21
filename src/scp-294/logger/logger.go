package logger

import (
	stlog "log"
	"os"
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

func Init(destination string) {
	log = stlog.New(fileLog(destination), "", stlog.LstdFlags)
}

func Log(message string) {
	log.Printf("%v\n", message)
}

func LogError(params ...interface{}) {
	log.Fatalf("[ERROR]%v\n", params)
}
