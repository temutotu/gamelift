package sessionhandler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gamelift-server-go/lib/packet"
	sessionmanager "gamelift-server-go/sessionsv/sessionManager"
)

func Leave(buf *[]byte) *[]byte {
	reader := bytes.NewReader(*buf)
	var packetLeave packet.LeaveRequest
	err := binary.Read(reader, binary.BigEndian, &packetLeave)
	if err != nil {
		fmt.Println("Read error:", err)
		return nil
	}
	fmt.Printf("playerID:%s", string(packetLeave.PlayerId[:]))

	sessionmanager.RemoveClient(string(packetLeave.PlayerId[:]))

	return nil
}
