package sessionmanager

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/amazon-gamelift/amazon-gamelift-servers-go-server-sdk/v5/server"
)

type SessionManager struct {
	SessionId  string
	Maplayer   int
	ClientList []SessionClient
}

type SessionClient struct {
	ID              string
	PlayerSessionID string
	IP              net.IP
	Port            int
	LastSeen        time.Time
	RTT             time.Duration
}

var inst *SessionManager

func GetSessionManager(sessionId string, maxPlayer int) *SessionManager {
	if inst == nil {
		inst = &SessionManager{
			SessionId: sessionId,
			Maplayer:  maxPlayer,
		}
	}
	return inst
}

func CheckClientJoin(playerId string) bool {
	self := GetSessionManager("", 0)
	for _, v := range self.ClientList {
		if v.ID == playerId {
			return true
		}
	}

	return false
}

func AcceptPlayer(playerId string, playerSessionId string, addr *net.UDPAddr) (int, error) {
	self := GetSessionManager("", 0)

	err := server.AcceptPlayerSession(playerSessionId)
	if err != nil {
		return 0, err
	}

	memberNo, err := self.AddClient(playerId, playerSessionId, addr)
	if err != nil {
		return 0, err
	}

	return memberNo, nil
}

func (self *SessionManager) AddClient(playerId string, playerSessionID string, addr *net.UDPAddr) (int, error) {
	memberNo := len(self.ClientList) + 1
	if len(self.ClientList) > self.Maplayer {
		return 0, errors.New("this session is full")
	}

	self.ClientList = append(self.ClientList, SessionClient{
		ID:              playerId,
		PlayerSessionID: playerSessionID,
		IP:              addr.IP,
		Port:            addr.Port,
		LastSeen:        time.Now(),
		RTT:             0,
	})
	fmt.Println("add player No.%d ID:%s", memberNo-1, playerId)

	return memberNo, nil
}

func RemoveClient(playerId string) {
	self := GetSessionManager("", 0)

	var removeIndex int
	var playerSessionID string
	for k, v := range self.ClientList {
		if v.ID == playerId {
			removeIndex = k
			playerSessionID = v.PlayerSessionID
			break
		}
	}

	err := server.RemovePlayerSession(playerSessionID)
	if err != nil {
		fmt.Println("Failed to RemovePlayerSession err:", err)
		return
	}

	self.ClientList = append(self.ClientList[:removeIndex], self.ClientList[removeIndex+1:]...)
}
