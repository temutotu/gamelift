package packet

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

type UnpackedData struct {
	IP       net.IP
	Port     int
	PlayerID string
	Payload  [236]byte
}

func PackPlayerData(playerIDStr string, payloadStr string) ([]byte, error) {
	// 2. UUIDのパース (16bytes)
	id, err := uuid.Parse(playerIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}

	// 3. 8バイト文字列の処理
	payloadBytes := []byte(payloadStr)
	if len(payloadBytes) > 252-16 {
		payloadBytes = payloadBytes[:252] // 切り捨て
	} else if len(payloadBytes) < 252 {
		// 8バイトに満たない場合はゼロパディング
		newBuf := make([]byte, 252)
		copy(newBuf, payloadBytes)
		payloadBytes = newBuf
	}

	// 4. バッファの作成
	buf := make([]byte, 252)

	copy(buf[0:16], id[:]) // UUIDをそのまま16byteコピー
	copy(buf[16:252], payloadBytes)

	return buf, nil
}

func UnpackPlayerData(data []byte) (*UnpackedData, error) {
	// 1. サイズチェック（予期しないパケットを弾く）
	if len(data) < 252 {
		return nil, fmt.Errorf("data too short: expected 28 bytes, got %d", len(data))
	}

	// 3. UUIDの復元 (4-20バイト)
	// バイナリ16バイトからUUIDオブジェクトを生成
	id, err := uuid.FromBytes(data[0:16])
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %v", err)
	}

	// 4. ペイロードの復元 (20-28バイト)
	// 文字列にキャスト（ゼロパディングが含まれる可能性がある点に注意
	payload := data[16:252]

	return &UnpackedData{
		PlayerID: id.String(),
		Payload:  [236]byte(payload),
	}, nil
}
