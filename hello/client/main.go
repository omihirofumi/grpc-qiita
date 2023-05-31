package main

import (
	pb "github.com/omihirofumi/grpc-go/hello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// 作成したgRPCサーバーのアドレス
var addr string = "localhost:50051"

func main() {
	// 作ったサーバーとのコネクタを作成
	// ssl認証なし
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("接続失敗: %v\n", err)
	}
	defer conn.Close()

	// gRPCクライアントの生成
	c := pb.NewHelloServiceClient(conn)

	doHello(c)
	doHelloAmp(c)
	doHelloManyTimes(c)
	doHelloEveryone(c)
}
