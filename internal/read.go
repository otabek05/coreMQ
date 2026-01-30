package internal

import (
	"errors"
	"io"
	"rawMQ/common"
)

func ReadPacket(r io.Reader) (*common.Packet, error) {
	// 1. Read first byte (packet type + flags)
	var header [1]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	packetType := header[0] >> 4
	flags := header[0] & 0x0F

	// 2. Read Remaining Length
	remainingLength, err := readRemainingLength(r)
	if err != nil {
		return nil, err
	}

	// 3. Read payload
	payload := make([]byte, remainingLength)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return &common.Packet{
		Type:    packetType,
		Flags:   flags,
		Length:  remainingLength,
		Payload: payload,
	}, nil
}



func readRemainingLength(r io.Reader) (int, error) {
	multiplier := 1
	value := 0

	for i := 0; i < 4; i++ {
		var encodedByte [1]byte
		if _, err := io.ReadFull(r, encodedByte[:]); err != nil {
			return 0, err
		}

		value += int(encodedByte[0]&127) * multiplier
		multiplier *= 128

		if encodedByte[0]&128 == 0 {
			return value, nil
		}
	}

	return 0, errors.New("malformed remaining length")
}

