package packet

import (
	"bytes"
	"encoding/binary"
	gameprotcol "gamelift-server-go/lib/gameProtcol"
)

type LeaveRequest struct {
	Type     uint8
	PlayerId [32]byte
}

func NewLeaveRequest(playerId string) ([]byte, error) {
	bytePlayerId := []byte(playerId)

	buf := new(bytes.Buffer)

	// 構造体の順に書き込んでいく
	if err := binary.Write(buf, binary.BigEndian, gameprotcol.PROTOCOL_LEAVE); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, bytePlayerId[:32]); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
