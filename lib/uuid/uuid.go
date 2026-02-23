package uuid

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUIDv4 は標準ライブラリのみでUUID v4を生成します
func GenerateUUIDv4() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// RFC 4122に準拠するためのビット操作
	// 1. バージョン4（0100）を設定
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// 2. バリアント（10xx）を設定
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// 8-4-4-4-12 の形式で文字列化
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}
