package sessionhandler

import (
	"fmt"
	gameprotcol "gamelift-server-go/lib/gameProtcol"
	"gamelift-server-go/lib/network"
	sessionmanager "gamelift-server-go/sessionsv/sessionManager"
	"net"
)

func Handle(buf *[]byte, resChan chan network.ResponseData, remoteAddr *net.UDPAddr) {
	protocol := (*buf)[0]
	var body *[]byte = nil

	if sessionmanager.CheckClientJoin(string((*buf)[1:33])) {
		switch protocol {
		case gameprotcol.PROTOCOL_LEAVE:
			body = Leave(buf)
		default:
			fmt.Println("no handler")
		}
		if body == nil {
			return
		}
	} else if protocol == gameprotcol.PROTOCOL_JOIN {
		body = Join(buf, remoteAddr)
	}

	if body == nil {
		return
	}

	resChan <- network.ResponseData{
		IP:   remoteAddr.IP,
		Port: remoteAddr.Port,
		Body: body,
	}
}
