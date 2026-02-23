package sessionhandler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gamelift-server-go/lib/packet"
	sessionmanager "gamelift-server-go/sessionsv/sessionManager"
	"net"
)

func Join(buf *[]byte, remoteAddr *net.UDPAddr) *[]byte {
	reader := bytes.NewReader(*buf)
	var packetJoin packet.Join
	err := binary.Read(reader, binary.BigEndian, &packetJoin)
	if err != nil {
		fmt.Println("Read error:", err)
		errMessage := "failed to read binary err:" + err.Error()
		return packet.NewErrorPakcet(errMessage)
	}
	fmt.Printf("playerID:%s,playerSessionID:%s", string(packetJoin.PlayerId[:]), string(packetJoin.PlayerSessionId[:]))

	memberNo, err := sessionmanager.AcceptPlayer(string(packetJoin.PlayerId[:]), string(packetJoin.PlayerSessionId[:]), remoteAddr)
	if err != nil {
		fmt.Println("Error Accept Player:", err)
		errMessage := "failed to accecpt player err:" + err.Error()
		return packet.NewErrorPakcet(errMessage)
	}

	return packet.NewJoinResponse(memberNo)
}
