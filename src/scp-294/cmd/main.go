package main

import (
	"github.com/edward/scp-294/internal"
	"github.com/edward/scp-294/internal/dbaccess"
	"github.com/edward/scp-294/internal/handler"
	_ "github.com/mattn/go-sqlite3"
	"myutil"
	"net/http"
)

func main() {
	internal.Logger = myutil.NewMyLogger("scp294.log")

	err := dbaccess.Connect()
	if err != nil {
		internal.Logger.Log("Main", err.Error())
		return
	}
	dbaccess.InitDatabases(false)

	server := http.Server{
		Addr: ":8294",
	}

	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/static/", fileServer)
	handler.RegisterRoutes()

	internal.Logger.Log("Main", "Server started!")
	server.ListenAndServe()
}
