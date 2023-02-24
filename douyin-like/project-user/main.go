package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//  初始化 db
	err := cli.InitDB()
	if err != nil {
		panic(err)
	}

	// 创建rpc服务
	rpcServer := grpc.NewServer()

	service.Dian - favoriteService(rpcServer, service.UserService)

	// 监听
	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal("启动监听失败！", err)
	}
	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("启动服务失败！", err)
	}
	fmt.Println("启动服务成功！")
}
