package broker

import (
	"log"
	"rawMQ/common"
)



func (b *Broker) handleSubscribe(client *common.Client, payload []byte) error {
	pos := 0

	packetID := int(payload[pos]) << 8 | int(payload[pos+1])
	pos += 2


	for pos < len(payload) {
		topic, err := b.readUTF8(payload, &pos)
		if err != nil {
			return err
		}

		qos := payload[pos]
		pos++

		log.Printf("SUBSCRIBE | client=%s topic=%s qos=%d", client.ID, topic, qos)

		b.addSubscription(topic, client)
	}

	b.sendSubAck(client.Conn, packetID)

	return nil 
}


func (b *Broker) addSubscription(topic string, client *common.Client) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscriptions[topic] = append(b.subscriptions[topic], client)
}
