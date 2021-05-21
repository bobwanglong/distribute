package main

import (
	"context"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
	"net/http"
)

func main() {
	host, port := "localhost", "30000"

	ctx, err := service.Start(context.Background(), "registryService",
		host, port, registerHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service")

}
func registerHandlers() {
	regService := registry.RegistryService{}
	http.HandleFunc("/services", regService.ServerHTTP)
}
