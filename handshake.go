package caddymc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var ErrTooBigVarint = errors.New("varint too large")
var ErrInvalidHandshake = errors.New("handshake is invalid")

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	Intent          int32
}

func ReadHandshakeFromByte(b []byte) (*Handshake, error) {
	return ReadHandshake(bytes.NewReader(b))
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	br := bufio.NewReader(r)

	length, err := ReadVarint(br)
	if err != nil {
		return nil, err
	}
	packet_id, err := ReadVarint(br)
	if err != nil {
		return nil, err
	}
	version, err := ReadVarint(br)
	if err != nil {
		return nil, err
	}
	hostname, err := ReadString(br)
	if err != nil {
		return nil, err
	}
	port, err := ReadShort(br)
	if err != nil {
		return nil, err
	}
	intent, err := ReadVarint(br)
	if err != nil {
		return nil, err
	}

	if packet_id != 0 {
		return nil, ErrInvalidHandshake
	}
	_ = length

	return &Handshake{
		ProtocolVersion: version,
		ServerAddress:   hostname,
		ServerPort:      port,
		Intent:          intent,
	}, nil
}

func ReadString(b io.ByteReader) (string, error) {
	out, err := ReadVarint(b)
	if err != nil {
		return "", err
	}

	return ReadStringLen(b, out)

}
func ReadStringLen(b io.ByteReader, ln int32) (string, error) {
	s := []rune{}
	for range ln {
		byte, err := b.ReadByte()
		if err != nil {
			return "", err
		}
		s = append(s, rune(byte))
	}
	return string(s), nil
}

func ReadVarint(b io.ByteReader) (int32, error) {
	var val int32 = 0
	for position := 0; position < 32; position += 7 {
		current_byte, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		val |= int32(current_byte&0x7F) << position

		if current_byte&0x80 == 0 {
			return val, nil
		}
	}

	return 0, ErrTooBigVarint
}

func ReadShort(b io.ByteReader) (uint16, error) {
	a1, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	a2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return 256*uint16(a1) + uint16(a2), nil
}
