package myutil

import (
	stlog "log"
	"os"
	"strconv"
)

var Logger *MyLogger

func NewMyLogger(logPath string) *MyLogger {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	log := stlog.New(file, "", stlog.LstdFlags)
	Logger = &MyLogger{
		file:   file,
		logger: log,
	}
	return Logger
}

type MyLogger struct {
	file   *os.File
	logger *stlog.Logger
}

func (my *MyLogger) Close() {
	my.file.Close()
}

func (my *MyLogger) Log(code string, message string) {
	my.logger.Printf("[%v] %v\n", code, message)
}

func (my *MyLogger) LogBytes(code string, message string, arr []byte, printDetails bool) {
	if printDetails {
		my.logger.Printf("[%v] %v\n%v\n", code, message, arr)
	} else {
		my.logger.Printf("[%v] %v\n", code, message)
	}
}

func (my *MyLogger) LogWarn(code string, message string) {
	my.logger.Printf("[%v] [WARN]%v\n", code, message)
}

func (my *MyLogger) LogError(code string, message string) {
	my.logger.Printf("[%v] [ERROR]%v\n", code, message)
}

func (my *MyLogger) GetCode(connId int, socket string) string {
	return "Conn" + strconv.Itoa(connId) + "-" + socket
}
