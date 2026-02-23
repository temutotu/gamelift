package main

import (
	"fmt"
	"gamelift-server-go/managesv/handler"
	"net/http"
)

func HTTPRouter() {
	// 1. ServeMux（ルーター）を生成
	mux := http.NewServeMux()

	// 2. ルーティングの定義
	// パスに "GET " などを付けると、特定のメソッドのみ許可できます（Go 1.22+）
	mux.HandleFunc("GET /Auth", handler.Authorize)
	mux.HandleFunc("POST /CreateGameSession", handler.CreateGameSession)
	mux.HandleFunc("POST /CreatePlayerSession", handler.CreatePlayerSession)
	mux.HandleFunc("POST /SearchGameSession", handler.SearchGameSession)
	// 3. サーバーの起動
	fmt.Println("Server starting on :8000...")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
