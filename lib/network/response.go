package network

import "net"

type ResponseData struct {
	IP   net.IP
	Port int
	Body *[]byte
}
