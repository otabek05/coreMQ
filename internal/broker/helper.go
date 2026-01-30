package broker

import (
	"io"
	"net"
	"rawMQ/common"
)




func (b *Broker) readUTF8(data []byte, pos *int) (string, error) {
	if *pos+2 > len(data) {
		return "", io.ErrUnexpectedEOF
	}

	length := int(data[*pos])<<8 | int(data[*pos+1])
	*pos += 2

	if *pos+length > len(data) {
		return "", io.ErrUnexpectedEOF
	}

	str := string(data[*pos : *pos+length])
	*pos += length
	return str, nil
}


func (b *Broker) sendSubAck(conn net.Conn, packetID int) error {
	packet := []byte{
		0x90, // SUBACK
		0x03, // remaining length
		byte(packetID >> 8),
		byte(packetID & 0xFF),
		0x00, // QoS 0 granted
	}
	_, err := conn.Write(packet)
	return err
}


func (b *Broker) handleConnect(payload []byte) (*common.ConnectInfo, error) {
	pos := 0

	protocolName, err := b.readUTF8(payload, &pos)
	if err != nil {
		return nil, err
	}

	protocolLevel := payload[pos]
	pos++

	flags := payload[pos]
	pos++

	keepAlive := int(payload[pos])<<8 | int(payload[pos+1])
	pos += 2

	clientID, err := b.readUTF8(payload, &pos)
	if err != nil {
		return nil, err
	}

	return &common.ConnectInfo{
		ProtocolName:  protocolName,
		ProtocolLevel: protocolLevel,
		ClientID:      clientID,
		KeepAlive:     keepAlive,
		Flags:         flags,
	}, nil
}


func (b *Broker) sendConnAck(conn net.Conn, returnCode byte) error {
	packet := []byte{
		0x20, // CONNACK
		0x02, // Remaining Length
		0x00, // Acknowledge Flags
		returnCode,
	}
	_, err := conn.Write(packet)
	return err
}


func (b *Broker) sendPingResp(conn net.Conn) error {
	packet := []byte{
		0xD0, // PINGRESP
		0x00, // Remaining Length
	}
	_, err := conn.Write(packet)
	return err
}
