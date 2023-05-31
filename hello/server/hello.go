package main

import (
	"context"
	"fmt"
	pb "github.com/omihirofumi/grpc-go/hello/proto" // protobufのインポート
	"io"
	"log"
	"strings"
)

// 1. 名前を送ると、”Hello FirstName”を返す
func (s *Server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Helloが呼び出されました: %v\n", in)
	// ここでは、文字列を生成し、レスポンスを返すだけ。
	return &pb.HelloResponse{
		Result: fmt.Sprintf("Hello %s", in.FirstName),
	}, nil
}

// 2. !が増幅して返ってくる
func (s *Server) HelloAmp(in *pb.HelloRequest, stream pb.HelloService_HelloAmpServer) error {
	log.Printf("HelloAmpが呼び出されました: %v\n", in)

	for i := 0; i < 5; i++ {
		// レスポンス作成
		res := fmt.Sprintf("Hello %s%s", in.FirstName, strings.Repeat("!", i))

		// クライアントにレスポンスを送信
		err := stream.Send(&pb.HelloResponse{Result: res})
		if err != nil {
			log.Fatalf("レスポンス送信中にエラーが発生: %v\n", err)
		}

	}

	return nil
}

// 3. 複数回送られてくるリクエストを結合して返す
func (s *Server) HelloManyTimes(stream pb.HelloService_HelloManyTimesServer) error {
	log.Println("HelloManyTimesが呼び出されました。")

	res := "Hello"

	// 無限ループでリクエストを受け取り続ける
	for {
		// リクエストを取得
		req, err := stream.Recv()

		// リクエスト終端の場合、レスポンスを送信
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{
				Result: res,
			})
		}

		// 終端以外のエラーの場合
		if err != nil {
			log.Fatalf("リクエスト読み取りでエラー発生: %v\n", err)
		}

		res = res + " " + req.FirstName
	}
}

// 4. リクエストごとに返信する
func (s *Server) HelloEveryone(stream pb.HelloService_HelloEveryoneServer) error {
	log.Println("HelloEveryoneが呼び出されました")

	for {
		// streamからリクエストを受け取り
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("リクエスト読み取り中にエラー発生: %v\n", err)
		}

		res := fmt.Sprintf("Hello %s", req.FirstName)
		// streamにレスポンスを送信
		err = stream.Send(&pb.HelloResponse{Result: res})
		if err != nil {
			log.Fatalf("レスポンス送信中にエラー発生: %v\n", err)
		}
	}
}
