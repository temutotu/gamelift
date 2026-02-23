package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gamelift-server-go/lib/packet"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type CreateGameSessionResPonse struct {
	PlayerSessionID   string
	SessionServerIP   string
	SessionServerPort int
}

func main() {
	var (
		uuid   = flag.String("u", "", "string flag")
		rating = flag.Int("r", 1000, "int flag")
	)

	flag.Parse()
	// シグナルを受け取るためのチャネルを作成
	sigs := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) と SIGTERM を監視対象にする
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	targetURL := "http://localhost:8000/SearchGameSession"

	values := url.Values{}
	values.Set("playerId", *uuid)
	values.Set("rating", strconv.Itoa(*rating))

	resp, err := http.PostForm(targetURL, values)
	if err != nil {
		fmt.Println("%v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	var res CreateGameSessionResPonse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("JSONのパースに失敗しました: %v", err)
		return
	}

	ip := res.SessionServerIP
	port := strconv.Itoa(res.SessionServerPort)

	joinReq, err := packet.NewJoinRequest(*uuid, res.PlayerSessionID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(ip + ":" + string(port))
	serverAddr, _ := net.ResolveUDPAddr("udp", ip+":"+string(port))
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()

	// サーバーへ送信
	conn.Write(joinReq)

	// サーバーからの返信を受信
	buf := make([]byte, 1024)
	go func() {
		sig := <-sigs // シグナルが来るまで待機
		fmt.Printf("\nシグナルを受信: %s", sig)
		leaveReq, err := packet.NewLeaveRequest(*uuid)
		if err != nil {
			return
		}

		conn.Write(leaveReq)
		os.Exit(0)
	}()
	for {
		n, _, _ := conn.ReadFromUDP(buf)
		fmt.Println("Reply from server: %#v\n", buf[:n])
	}
}
