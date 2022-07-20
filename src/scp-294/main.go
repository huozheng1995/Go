package main

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/controller"
	"github.com/edward/scp-294/logger"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func init() {
	logger.Init("scp294.log")
	common.Connect()
	common.InitDatabases()
}

func main() {
	server := http.Server{
		Addr: common.HttpAddr,
	}

	fileServer := http.FileServer(http.Dir("wwwroot"))
	http.Handle("/js/", fileServer)
	http.Handle("/css/", fileServer)
	http.Handle("/img/", fileServer)
	controller.RegisterRoutes()

	logger.Log("Server starting...")
	server.ListenAndServe()
}
