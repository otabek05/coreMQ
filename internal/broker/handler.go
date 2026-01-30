package broker

import (
	"log"
	"net"
	"rawMQ/common"
	"rawMQ/internal"
)

const (
	ConnAccepted          = 0x00
	ConnRefusedIdentifier = 0x02
)

func (b *Broker) HandleConn(conn net.Conn) {
	defer conn.Close()

	var clientID string

	for {
		packet, err := internal.ReadPacket(conn)
		if err != nil {
			log.Println("connection closed:", err)
			if clientID != "" {
				b.removeClient(clientID)
			}
			return
		}

		switch packet.Type {

		case 1: // CONNECT
			info, err := b.handleConnect(packet.Payload)
			if err != nil {
				b.sendConnAck(conn, ConnRefusedIdentifier)
				return
			}

			clientID = info.ClientID

			client := &common.Client{
				ID:   clientID,
				Conn: conn,
			}

			if _, exists := b.getClient(clientID); exists {
				log.Println("Duplicate client ID, disconnecting old one:", clientID)
				b.removeClient(clientID)
			}

			b.addClient(client)
			b.sendConnAck(conn, ConnAccepted)

		case 8: // SUBSCRIBE
			client, ok := b.getClient(clientID)
			if !ok {
				log.Println("SUBSCRIBE from unknown client")
				return
			}
			b.handleSubscribe(client, packet.Payload)

		case 3: // PUBLISH
			client, ok := b.getClient(clientID)
			if !ok {
				log.Println("PUBLISH from unknown client")
				return
			}
			b.handlePublish(packet, client)
		
		case 12:
			log.Println("Ping message received from client", clientID)
			b.sendPingResp(conn)

		case 14: // DISCONNECT
			log.Println("Client disconnected:", clientID)
			b.removeClient(clientID)
			return
		}
	}
}

