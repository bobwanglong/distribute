package main

import (
	"context"
	"distributed/log"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./cmd/logservice/distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(context.Background(), "Log Service", host,
		port, log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done() //当调用cancel函数的时候ctx就有Done

	fmt.Println("Shutting down log service")
}
