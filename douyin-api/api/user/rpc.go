package user

import (
	"github.com/codeRookieC/douyin/douyin-grpc/user/userServer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var UserServerClient userServer.UserServiceClient

func InitRpcUserClient() {
	conn, err := grpc.Dial(":1001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	UserServerClient = userServer.NewUserServiceClient(conn)
}
