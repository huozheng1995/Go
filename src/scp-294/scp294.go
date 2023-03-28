package main

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/controller"
	"github.com/edward/scp-294/logger"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	logger.Init("scp294.log")
	err := common.Connect()
	if err != nil {
		logger.Log(err.Error())
		return
	}
	common.InitDatabases(true)

	server := http.Server{
		Addr: common.HttpAddr,
	}

	fileServer := http.FileServer(http.Dir("wwwroot"))
	http.Handle("/js/", fileServer)
	http.Handle("/css/", fileServer)
	http.Handle("/img/", fileServer)
	controller.RegisterRoutes()

	logger.Log("Server started!")
	server.ListenAndServe()
}
