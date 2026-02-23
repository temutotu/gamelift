package main

import (
	"flag"
	"fmt"
	gameprotcol "gamelift-server-go/lib/gameProtcol"
	"gamelift-server-go/lib/network"
	sessionhandler "gamelift-server-go/sessionsv/sessionHandler"
	sessionmanager "gamelift-server-go/sessionsv/sessionManager"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/amazon-gamelift/amazon-gamelift-servers-go-server-sdk/v5/model"
	"github.com/amazon-gamelift/amazon-gamelift-servers-go-server-sdk/v5/server"
)

func main() {
	var (
		processID = flag.String("i", "process-123", "string flag")
		port      = flag.Int("p", 7777, "int flag")
	)

	flag.Parse()

	url := "http://localhost:8000/Auth"

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get AuthToken %v\n", err)
		return
	}
	defer res.Body.Close()

	// 6. レスポンスボディの読み取り
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Println(string(responseBody))
	// GameLift Local (9334ポート) への接続設定
	// SDK v5では GameLift Local も WebSocket プロトコルで通信します
	params := server.ServerParameters{
		WebSocketURL: "wss://us-east-1.api.amazongamelift.com",     // Local jarの場合はここを指定
		ProcessID:    *processID,                                   // 任意のID
		HostID:       "local",                                      // 任意のID
		FleetID:      "fleet-55b311f9-dc09-4454-a704-d4855789e109", // 任意のID
		AuthToken:    string(responseBody),
	}

	err = server.InitSDK(params)
	if err != nil {
		log.Fatalf("SDK初期化失敗: %v", err)
	}

	done := make(chan struct{})
	// 2. ハンドラーの実装
	// GameLiftサービスからサーバープロセスへの各種通知を処理するコールバック
	processParameters := &server.ProcessParameters{
		// ゲームセッションが割り当てられた時に呼ばれる
		OnStartGameSession: func(gameSession model.GameSession) {
			log.Printf("ゲームセッション開始: %s", gameSession.GameSessionID)
			// セッションをアクティブにする（これでプレイヤーが入室可能になる）
			server.ActivateGameSession()
			done <- struct{}{}
		},
		// プロセスが終了を求められた時に呼ばれる
		OnProcessTerminate: func() {
			log.Println("プロセス終了通知を受信")
			server.ProcessEnding()
			os.Exit(0)
		},
		// GameLiftへの死活監視（1分間に数回呼ばれる）
		OnHealthCheck: func() bool {
			return true // サーバーが正常ならtrue
		},
		// ログファイルのパスを指定
		LogParameters: server.LogParameters{
			LogPaths: []string{"/local/game/logs/myserver.log"},
		},
		// このプロセスがリッスンするポート
		Port: *port,
	}
	fmt.Println("GameLiftサーバが起動開始")

	// 3. GameLiftサービスに「準備完了」を通知
	err = server.ProcessReady(*processParameters)
	if err != nil {
		log.Fatalf("ProcessReady失敗: %v", err)
	}

	// サーバープロセスを維持（実際にはここでゲームループやgRPC/WebSocketサーバなどを動かす）
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	fmt.Println("GameLiftサーバが起動し、待機状態になりました Port:" + strconv.Itoa(*port))
	<-done
	sessionId, err := server.GetGameSessionID()
	if err != nil {
		log.Fatalf("failed to get sessionID: %v", err)
	}

	sessionmanager.GetSessionManager(sessionId, 4)
	fmt.Println("Created SessionManager")

	// 2. UDPコネクションの作成
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on :8080...")

	resChan := make(chan network.ResponseData)
	go func() {
		buffer := make([]byte, 1024)

		for {
			// 3. データの受信
			// n: 受信バイト数, remoteAddr: 送信元のIP/ポート
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}

			message := string(buffer[:n])
			fmt.Printf("Received from %v: %s\n", remoteAddr, message)

			protocolNo := buffer[0]
			fmt.Println("Protocol is", protocolNo)

			if protocolNo > gameprotcol.PROTOCOL_ERROR {
				// game session process
			} else {
				sessionhandler.Handle(&buffer, resChan, remoteAddr)
			}
		}
	}()

	for {
		resData := <-resChan

		resAddr := &net.UDPAddr{
			IP:   resData.IP,
			Port: resData.Port,
		}

		conn.WriteToUDP(*resData.Body, resAddr)
	}
}
