package main

import (
	"context"
	"fmt"
	"github.com/edward/distributed/log"
	"github.com/edward/distributed/registry"
	"github.com/edward/distributed/service"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down log service.")
}
