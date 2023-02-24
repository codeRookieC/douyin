package main

import (
	"context"
	"fmt"
	"log"
	"video_server/api/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer conn.Close()

	Client := pb.NewVideoServiceClient(conn)

	// ........... testpb.DouyinPublishAction..........
	// var token string
	// token = "my token"
	// var title string
	// title = "title"
	// file, err := os.Open("../resource/go.mp4")
	// if err != nil {
	// 	fmt.Println("Open file err =", err)
	// 	return
	// }
	// defer file.Close()
	// var data []byte
	// buf := make([]byte, 1024)
	// for {
	// 	n, err := file.Read(buf)
	// 	if err != nil && err != io.EOF {
	// 		fmt.Println("read buf fail", err)
	// 		return
	// 	}
	// 	//说明读取结束
	// 	if n == 0 {
	// 		break
	// 	}
	// 	//读取到最终的缓冲区中
	// 	data = append(data, buf[:n]...)
	// }
	// request := pb.DouyinPublishActionRequest{
	// 	Token: &token,
	// 	Data:  data,
	// 	Title: &title,
	// }
	// response, err := Client.DouyinPublishAction(context.Background(), &request)
	// if err != nil {
	// 	log.Fatal("DouyinFeed error", err)
	// }
	// fmt.Print(response)
	//
	// ...................Client.DouyinFeed(context.Background(), &request)..............................
	// var time int64
	// var token string
	// token = "my token"

	// time = 0
	// request := pb.DouyinFeedRequest{
	// 	LatestTime: &time,
	// 	Token:      &token,
	// }
	// response, err := Client.DouyinFeed(context.Background(), &request)
	// if err != nil {
	// 	log.Fatal("DouyinFeed error", err)
	// }
	// fmt.Println(response)

	// ..................Client.DouyinPublishList(context.Background(), &request)......................

	var token string
	token = "my token"
	var UserId int64
	UserId = 001
	request := pb.DouyinPublishListRequest{
		UserId: &UserId,
		Token:  &token,
	}
	response, err := Client.DouyinPublishList(context.Background(), &request)
	if err != nil {
		log.Fatal("DouyinFeed error", err)
	}
	fmt.Print(response)

}
