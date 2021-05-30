package main

import (
	"context"
	"distributed/grades"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main(){
	host,port := "localhost","6000"
	serviceAddress :=fmt.Sprintf("http://%s:%s",host,port)
	
	r:=registry.Registration{
		ServiceName: registry.GradingService,
		ServiceUrl: serviceAddress,
	}
	ctx,err :=service.Start(context.Background(),r,host,port,grades.RegisterHandlers)
	if err != nil{
		stlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down grading service")

}