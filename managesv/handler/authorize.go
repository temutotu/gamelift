package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	// 1. AWS設定のロード (環境変数や共有設定ファイルから読み込み)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // 適切なリージョンを指定
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// 2. GameLiftクライアントの作成
	client := gamelift.NewFromConfig(cfg)

	// 3. 入力パラメータの設定
	fleetID := "fleet-55b311f9-dc09-4454-a704-d4855789e109"
	computeName := "local"

	input := &gamelift.GetComputeAuthTokenInput{
		FleetId:     &fleetID,
		ComputeName: &computeName,
	}

	// 4. APIの実行
	result, err := client.GetComputeAuthToken(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get compute auth token, %v", err)
	}

	fmt.Fprintf(w, aws.ToString(result.AuthToken))
}
