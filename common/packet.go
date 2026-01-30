package common

import (
	"net"
	"sync"
)


type Client struct {
	ID   string
	Conn net.Conn
}

type Broker struct {
	mu            sync.RWMutex
	subscriptions map[string][]*Client // topic -> clients
}


type Packet struct {
	Type     byte
	Flags    byte
	Length   int
	Payload  []byte
}



type ConnectInfo struct {
	ProtocolName  string
	ProtocolLevel byte
	ClientID      string
	KeepAlive     int
	Flags         byte
}