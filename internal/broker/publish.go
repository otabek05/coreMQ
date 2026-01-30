package broker

import (
	"encoding/binary"
	"log"
	"net"
	"rawMQ/common"
)


func (b *Broker) handlePublish( packet *common.Packet, client *common.Client) error {
	pos := 0
	payload := packet.Payload

	topic, err := b.readUTF8(payload, &pos)
	if err != nil {
		return err
	}

	message := payload[pos:]

	log.Printf(
		"PUBLISH | from=%s topic=%s payload=%s",
		client.ID,
		topic,
		string(message),
	)

	b.publish(topic, message)
	return nil
}


func (b *Broker) publish(topic string, message []byte) {
	b.mu.RLock()
	subs := b.subscriptions[topic]
	b.mu.RUnlock()

	for _, sub := range subs {
		sendPublish(sub.Conn, topic, message)
	}
}


func sendPublish(conn net.Conn, topic string, payload []byte) error {
	var body []byte
	body = appendString(body, topic)
	body = append(body, payload...)

	packet := []byte{0x30} // PUBLISH QoS 0
	packet = append(packet, encodeRemainingLength(len(body))...)
	packet = append(packet, body...)

	_, err := conn.Write(packet)
	return err
}


func appendString(buf []byte, s string) []byte {
	var lenBytes [2]byte
	binary.BigEndian.PutUint16(lenBytes[:], uint16(len(s)))

	buf = append(buf, lenBytes[:]...)
	buf = append(buf, []byte(s)...)
	return buf
}

func encodeRemainingLength(length int) []byte {
	var encoded []byte

	for {
		digit := length % 128
		length /= 128

		// if there are more digits, set the continuation bit
		if length > 0 {
			digit |= 0x80
		}

		encoded = append(encoded, byte(digit))

		if length == 0 {
			break
		}
	}

	return encoded
}
