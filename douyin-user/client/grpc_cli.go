package main

import (
	"context"
	"fmt"
	"github.com/codeRookieC/douyin/douyin-user/client/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// 1. 新建连接，端口是服务端开放的8002端口
	// 没有证书会报错
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// 退出时关闭链接
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// 2. 调用Product.pb.go中的NewProdServiceClient方法
	userServiceClient := service.NewUserServiceClient(conn)

	// 3. 直接像调用本地方法一样调用GetProductStock方法
	resp, err := userServiceClient.Login(context.Background(), &service.DouyinUserLoginRequest{Username: "ab", Password: "a"})
	if err != nil {
		log.Fatal("调用gRPC方法错误: ", err)
	}

	fmt.Println("调用gRPC方法成功", resp.UserId)
}
