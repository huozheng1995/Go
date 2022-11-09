package main

import (
	"context"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/controller"
	"github.com/edward/scp-294/logger"
	"github.com/kardianos/service"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	logger.Init("scp294.log")
	srvConfig := &service.Config{
		Name:        "SCP294",
		DisplayName: "SCP294",
		Description: "SCP294 data converter",
	}
	prg := &program{}
	s, err := service.New(prg, srvConfig)
	if err != nil {
		logger.LogError(err.Error())
	}

	err = s.Run()
	if err != nil {
		logger.LogError("Failed to run SCP294: ", err)
	}
}

type program struct{}

var server http.Server
var ctx context.Context

func (p *program) Start(srv service.Service) error {
	logger.Log("Starting SCP294 server...")
	ctx = context.Background()
	go p.run()
	srv.Start()
	return nil
}

func (p *program) run() {
	err := common.Connect()
	if err != nil {
		logger.Log(err.Error())
	}
	common.InitDatabases(false)

	server = http.Server{
		Addr: common.HttpAddr,
	}

	fileServer := http.FileServer(http.Dir("wwwroot"))
	http.Handle("/js/", fileServer)
	http.Handle("/css/", fileServer)
	http.Handle("/img/", fileServer)
	controller.RegisterRoutes()

	logger.Log("SCP294 server is started!")
	logger.LogError(server.ListenAndServe())
}

func (p *program) Stop(srv service.Service) error {
	logger.Log("Stopping SCP294 server...")
	server.Shutdown(ctx)
	srv.Stop()
	return nil
}
