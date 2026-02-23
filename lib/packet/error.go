package packet

import (
	gameprotcol "gamelift-server-go/lib/gameProtcol"
)

type Error struct {
	Type    uint8
	Message string
}

func NewErrorPakcet(message string) *[]byte {
	buf := []byte(message)
	length := uint16(len(buf))

	buffer := make([]byte, 1+2+len(buf))
	buffer[0] = gameprotcol.PROTOCOL_ERROR
	buffer[1] = byte(length >> 8)
	buffer[2] = byte(length & 0xFF)
	copy(buffer[2:], buf)

	return &buffer
}
