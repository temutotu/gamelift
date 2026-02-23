package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gamelift-server-go/lib/network"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

func CreatePlayerSession(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析エラー", http.StatusBadRequest)
		return
	}

	gameSessionId := r.PostFormValue("gameSessionId")
	playerId := r.PostFormValue("playerId")
	if gameSessionId == "" || playerId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "param gameSessionId or playerId is empty")
		return
	}
	// AWS設定のロード
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "AWS設定ロード失敗: %w", err)
		return
	}

	client := gamelift.NewFromConfig(cfg)

	// セッション作成のパラメータ設定
	input := &gamelift.CreatePlayerSessionInput{
		GameSessionId: aws.String(gameSessionId),
		PlayerId:      aws.String(playerId),
	}

	// 実行
	output, err := client.CreatePlayerSession(context.TODO(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "プレイヤーの参加リクエスト失敗: %v", err)
		return
	}

	response := network.ReseponseManageServer{
		PlayerSessionID:   *output.PlayerSession.PlayerSessionId,
		SessionServerIP:   *output.PlayerSession.IpAddress,
		SessionServerPort: *output.PlayerSession.Port,
	}

	w.WriteHeader(http.StatusCreated)

	// 4. 構造体をJSONに変換してResponseWriterに書き込む
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// エンコードに失敗した場合のエラーハンドリング
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
