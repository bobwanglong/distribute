package main

import (
	"context"
	"distributed/grades"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main(){
	host,port := "localhost","6001"
	serviceAddress :=fmt.Sprintf("http://%s:%s",host,port)
	
	r:=registry.Registration{
		ServiceName: registry.GradingService,
		ServiceUrl: serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress+"/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx,err :=service.Start(context.Background(),r,host,port,grades.RegisterHandlers)
	if err != nil{
		stlog.Fatal(err)
	}
	if logProvider,err :=registry.GetProvider(registry.LogService); err ==nil{
		fmt.Println("Loggin service found at:",logProvider)
		log.SetClientLogger(logProvider,r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("shutting down grading service")

}