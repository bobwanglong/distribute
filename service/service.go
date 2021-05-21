package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, serviceName, host, port string, registerHandlersfunc func()) (context.Context, error) {
	registerHandlersfunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe()) // 如果服务启动发生错误就返回该错误
		cancel()                          //调用上下文的cancel关闭
	}()
	//用户可以手动停止服务
	go func() {
		fmt.Printf("%s started,press any key to stop\n", serviceName)
		// 只要用户按任何键就会往下执行，否则就停在此处
		var s string
		fmt.Scanln(&s)
		//继续执行就停止服务
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
