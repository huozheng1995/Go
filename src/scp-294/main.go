package main

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/controller"
	"github.com/edward/scp-294/logger"
	_ "github.com/mattn/go-sqlite3"
	"myutil"
	"net/http"
)

func main() {
	logger.Logger = myutil.NewMyLogger("scp294.log")

	err := common.Connect()
	if err != nil {
		logger.Logger.Log("Main", err.Error())
		return
	}
	common.InitDatabases(false)

	server := http.Server{
		Addr: ":8294",
	}

	fileServer := http.FileServer(http.Dir("wwwroot"))
	http.Handle("/js/", fileServer)
	http.Handle("/css/", fileServer)
	http.Handle("/img/", fileServer)
	controller.RegisterRoutes()

	logger.Logger.Log("Main", "Server started!")
	server.ListenAndServe()
}
