package handler

import (
	"context"
	"errors"
	"fmt"
	"gamelift-server-go/consts"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/smithy-go"
)

func SearchGameSession(w http.ResponseWriter, r *http.Request) {
	fleetID := "fleet-55b311f9-dc09-4454-a704-d4855789e109"

	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析エラー", http.StatusBadRequest)
		return
	}

	playerId := r.PostFormValue("playerId")
	if playerId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "param playerId is empty")
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

	// セッション作成のパラメータ設定
	fileterExpression := fmt.Sprintf("gameSessionProperties.rank = '%s'", rank)
	input := &gamelift.SearchGameSessionsInput{
		FleetId:          aws.String(fleetID),
		FilterExpression: aws.String(fileterExpression),
		SortExpression:   aws.String("creationTimeMillis ASC"),
		Limit:            aws.Int32(10),
	}

	// 実行
	output, err := client.SearchGameSessions(context.TODO(), input)
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			// 具体的なエラーコードやメッセージを表示
			fmt.Printf("Operation Error: %s, %s, %v\n", oe.Service(), oe.Operation(), oe.Err)
		} else {
			fmt.Println("Generic Error:", err)
		}
	}

	if err != nil || len(output.GameSessions) == 0 {
		r.URL.Path = "/CreateGameSession"
		CreateGameSession(w, r)
	} else {
		r.PostForm.Set("gameSessionId", *output.GameSessions[0].GameSessionId)
		CreatePlayerSession(w, r)
	}
}
