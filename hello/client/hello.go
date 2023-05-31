package main

import (
	"context"
	pb "github.com/omihirofumi/grpc-go/hello/proto"
	"io"
	"log"
	"time"
)

// 1. Hello呼び出し
func doHello(c pb.HelloServiceClient) {
	log.Println("---doHello---")

	// リクエストを作成
	req := &pb.HelloRequest{FirstName: "Hokke"}

	// Helloサービス呼び出し
	res, err := c.Hello(context.Background(), req)
	if err != nil {
		log.Fatalf("Helloの呼び出しでエラー: %v\n", err)
	}
	log.Printf("%s\n", res.Result)
}

// 2. HelloAmp呼び出し
func doHelloAmp(c pb.HelloServiceClient) {
	log.Println("---doHelloAmp---")

	// リクエスト作成
	req := &pb.HelloRequest{FirstName: "Hokke"}

	// gRPCとのstream生成
	stream, err := c.HelloAmp(context.Background(), req)
	if err != nil {
		log.Fatalf("HelloAmp呼び出しでエラー: %v\n", err)
	}

	// EOFが送られるまでレスポンスを取得し続けるため、無限ループ
	for {
		// レスポンス取得
		res, err := stream.Recv()

		// レスポンスが終わった場合
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("レスポンス取得でエラー発生: %v\n", err)
		}

		log.Printf("Response is %s\n", res.Result)
	}
}

// 3. HelloManyTimes呼び出し
func doHelloManyTimes(c pb.HelloServiceClient) {
	log.Println("---doHelloManyTimes---")

	// stream生成
	stream, err := c.HelloManyTimes(context.Background())
	if err != nil {
		log.Fatalf("HelloManyTimes呼び出しでエラー")
	}

	// リクエスト生成
	reqs := []*pb.HelloRequest{
		{FirstName: "Hokke"},
		{FirstName: "Hiro"},
		{FirstName: "Hoge"},
	}

	// リクエスト数分送信
	for _, req := range reqs {
		log.Printf("Sending Request: %v\n", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("レスポンス取得でエラー: %v\n", err)
	}

	log.Println(res.Result)
}

// 4. HelloEveryone呼び出し
func doHelloEveryone(c pb.HelloServiceClient) {
	log.Println("---doHelloEveryone---")

	stream, err := c.HelloEveryone(context.Background())
	if err != nil {
		log.Fatalf("呼び出しでエラー")
	}

	reqs := []*pb.HelloRequest{
		{FirstName: "Hokke"},
		{FirstName: "Hiro"},
		{FirstName: "Hoge"},
	}

	// ゴルーチンを使うので、待機チャネルを使用
	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Sending request: %v\n", req)
			stream.Send(req)
			// わかりやすいように、１s待つ
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("レスポンス取得中にエラー: %v\n", err)
			}
			log.Printf("Response is %s\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc
}
