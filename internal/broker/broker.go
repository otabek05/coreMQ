package broker

import (
	"rawMQ/common"
	"sync"
)


type Broker struct {
	mu            sync.RWMutex
	clients  map[string]*common.Client
	subscriptions map[string][]*common.Client // topic -> clients
}


func New() *Broker {
	return &Broker{
		clients: make(map[string]*common.Client),
		subscriptions: make(map[string][]*common.Client),
	}
}


