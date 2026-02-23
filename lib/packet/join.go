package packet

import (
	"bytes"
	"encoding/binary"
	gameprotcol "gamelift-server-go/lib/gameProtcol"
)

type Join struct {
	Type            uint8
	PlayerId        [32]byte
	PlayerSessionId [42]byte
}

type JoinResponse struct {
	Type     uint8
	MemberNo uint8
}

func NewJoinRequest(playerId string, playerSessionId string) ([]byte, error) {
	bytePlayerId := []byte(playerId)
	bytePlayerSessionId := []byte(playerSessionId)

	buf := new(bytes.Buffer)

	// 構造体の順に書き込んでいく
	if err := binary.Write(buf, binary.BigEndian, gameprotcol.PROTOCOL_JOIN); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, bytePlayerId[:32]); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, bytePlayerSessionId[:42]); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewJoinResponse(memberNo int) *[]byte {
	buffer := make([]byte, 2)
	buffer[0] = gameprotcol.PROTOCOL_JOIN
	buffer[1] = byte(memberNo)

	return &buffer
}
