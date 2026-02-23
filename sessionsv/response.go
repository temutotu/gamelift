package main

import "net"

type responseData struct {
	IP   net.IP
	Port int
	Body *[]byte
}
