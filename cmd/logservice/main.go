package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./cmd/logservice/distributed.log")
	host, port := "localhost", "4000"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: registry.LogService,
		ServiceUrl:  serviceAddr,
	}
	ctx, err := service.Start(context.Background(), r, host,
		port, log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done() //当调用cancel函数的时候ctx就有Done

	fmt.Println("Shutting down log service")
}
