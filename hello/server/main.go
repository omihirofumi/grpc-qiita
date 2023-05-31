package main

import (
	pb "github.com/omihirofumi/grpc-go/hello/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

// gRPCサーバーで受け付けるアドレス
var addr string = "0.0.0.0:50051"

// gRPCサーバーの構造体を宣言
type Server struct {
	pb.HelloServiceServer
}

func main() {
	// 指定したアドレスでTCPプロトコルで通信受付インスタンス生成
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listen on %s\n", addr)

	// grpcサーバーを生成し、Helloサービスを設定
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &Server{})

	// サーバー起動
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
