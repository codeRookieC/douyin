package main

import (
	"log"
	"net"
	"video_server/api/proto/pb"
	"video_server/config"
	"video_server/service"

	"google.golang.org/grpc"
)

func main() {
	rpcServer := grpc.NewServer()
	pb.RegisterVideoServiceServer(rpcServer, service.New())
	listener, err := net.Listen("tcp", config.C.GC.Addr)
	if err != nil {
		log.Fatal("listen error", err)
	}
	err = rpcServer.Serve(listener)

}
