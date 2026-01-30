package main

import (
	"log"
	"net"
	"rawMQ/internal/broker"
)


func main() {
	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		panic(err)
	}

	log.Println("MQTT Broker listening on port :1883")
	broker := broker.New()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go broker.HandleConn(conn)

	}
}

