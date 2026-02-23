package handler

import (
	"context"
	"fmt"
	"gamelift-server-go/consts"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

func CreateGameSession(w http.ResponseWriter, r *http.Request) {
	fleetID := "fleet-55b311f9-dc09-4454-a704-d4855789e109"
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析エラー", http.StatusBadRequest)
		return
	}

	ratingStr := r.PostFormValue("rating")
	if ratingStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ratingStr is empty")
		return
	}
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "rating is not numeric")
		return
	}

	rank := consts.GetRankLabel(rating)

	// AWS設定のロード
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "AWS設定ロード失敗: %w", err)
		return
	}

	client := gamelift.NewFromConfig(cfg)

	fmt.Println(rank)
	// セッション作成のパラメータ設定
	input := &gamelift.CreateGameSessionInput{
		FleetId:                   aws.String(fleetID),
		MaximumPlayerSessionCount: aws.Int32(4), // 最大プレイヤー数
		Location:                  aws.String("custom-local-test"),
		GameProperties: []types.GameProperty{
			{
				Key:   aws.String("rank"),
				Value: aws.String(rank),
			},
		},
	}

	// 実行
	output, err := client.CreateGameSession(context.TODO(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ゲームセッション作成失敗: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}
	r.PostForm.Set("gameSessionId", *output.GameSession.GameSessionId)
	CreatePlayerSession(w, r)
}
