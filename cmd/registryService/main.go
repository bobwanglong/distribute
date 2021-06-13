package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

// import (
// 	"context"
// 	"distributed/registry"
// 	"distributed/service"
// 	"fmt"
// 	stlog "log"
// 	"net/http"
// )

// func main() {
// 	host, port := "localhost", "30000"

// 	ctx, err := service.Start(context.Background(), "registryService",
// 		host, port, registerHandlers)
// 	if err != nil {
// 		stlog.Fatalln(err)
// 	}
// 	<-ctx.Done()
// 	fmt.Println("Shutting down log service")

// }
// func registerHandlers() {
// 	regService := registry.RegistryService{}
// 	http.HandleFunc("/services", regService.ServerHTTP)
// }

func main(){
	//心跳检测
	registry.SetupRegistryService()
	
	http.Handle("/services",&registry.RegistryService{})
	ctx,cancel :=context.WithCancel(context.Background())
	defer cancel()
	var srv http.Server
	srv.Addr=registry.ServerPort
	
	go func ()  {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func ()  {
		fmt.Println("RegisterServer started,press any key to stop")
		// 只要用户按任何键就会往下执行，否则就停在此处
		var s string
		fmt.Scanln(&s)
		//继续执行就停止服务
		srv.Shutdown(ctx)
		cancel()	
	}()
	<-ctx.Done()
		fmt.Println("Shutting down Registry service")
}