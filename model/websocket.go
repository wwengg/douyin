package model

import "encoding/binary"

type websocketPacket struct {
	Valid         bool
	flags         byte
	opcode        int
	opcode_str    string
	mask          bool
	payloadLength int
	maskingKey    []byte
	Payload       []byte
	PacketSize    int
}

func NewWebsocketPacket(packet []byte) *websocketPacket {
	p := &websocketPacket{}
	if len(packet) <= 6 {
		p.Valid = false
		return p
	}
	p.Valid = true
	p.flags = packet[0] & 0xF0
	p.opcode = int(packet[0] & 0x0F)
	p.mask = (packet[1] & 0x80) == 0x80
	p.payloadLength = int(packet[1] & 0x7f)
	packetStart := 2
	if p.payloadLength == 126 {
		if len(packet) <= 8 {
			p.Valid = false
			return p
		}

		p.payloadLength = int(packet[2])<<8 | int(packet[3])
		packetStart += 2
		if p.mask {
			p.maskingKey = packet[4:8]
		}
	} else if p.payloadLength == 127 {
		if len(packet) <= 14 {
			p.Valid = false
			return p
		}

		p.payloadLength = int(packet[2])<<56 | int(packet[3])<<48 | int(packet[4])<<40 | int(packet[5])<<32 | int(packet[6])<<24 | int(packet[7])<<16 | int(packet[8])<<8 | int(packet[9])
		packetStart += 8
		if p.mask {
			p.maskingKey = packet[10:14]
		}
	} else {
		if p.mask {
			p.maskingKey = packet[2:6]
		}
	}

	if p.mask {
		packetStart += 4
	}

	p.PacketSize = packetStart + p.payloadLength
	if packetStart > len(packet) || p.PacketSize > len(packet) {
		p.Payload = make([]byte, 0)
		p.Valid = false
	} else {
		p.Payload = packet[packetStart:p.PacketSize]
		p.Valid = (len(p.Payload) == p.payloadLength)
	}

	if p.opcode == 0x00 {
		p.opcode_str = "Continuation"
	} else if p.opcode == 0x01 {
		p.opcode_str = "Text"
	} else if p.opcode == 0x02 {
		p.opcode_str = "Binary"
	} else if p.opcode == 0x08 {
		p.opcode_str = "Close"
	} else if p.opcode == 0x09 {
		p.opcode_str = "Ping"
	} else if p.opcode == 0x0A {
		p.opcode_str = "Pong"
	} else {
		p.opcode_str = "Unknown"
	}

	if p.Valid && p.maskingKey != nil {
		for i := 0; i < p.payloadLength; i++ {
			p.Payload[i] ^= p.maskingKey[i%4]
		}
	}

	return p
}

func (packet *websocketPacket) Encode() []byte {
	packetLength := 2 + len(packet.maskingKey) + len(packet.Payload)
	packetStart := 2 + len(packet.maskingKey)
	maskingKeyStart := 2

	packet.payloadLength = len(packet.Payload)

	if packet.payloadLength > 125 {
		packetLength += 2
		packetStart += 2
		maskingKeyStart += 2
	}

	if packet.payloadLength > 65535 {
		packetLength += 6
		packetStart += 6
		maskingKeyStart += 6
	}

	buf := make([]byte, packetLength)
	buf[0] = byte(packet.flags) | byte(packet.opcode)

	// encode the length
	if packet.payloadLength < 126 {
		buf[1] = byte(packet.payloadLength)
	} else if packet.payloadLength < 65536 {
		buf[1] = 126
		binary.BigEndian.PutUint16(buf[2:], uint16(packet.payloadLength))
	} else {
		buf[1] = 127
		binary.BigEndian.PutUint64(buf[2:], uint64(packet.payloadLength))
	}

	if packet.maskingKey != nil {
		buf[1] |= 0x80
	}

	// reencode using the masking key
	if packet.maskingKey != nil {
		copy(buf[maskingKeyStart:], packet.maskingKey)
		copy(buf[packetStart:], packet.Payload)
		for i := 0; i < len(packet.Payload); i++ {
			buf[i+packetStart] ^= packet.maskingKey[i%4]
		}
	} else {
		copy(buf[packetStart:], packet.Payload)
	}

	return buf
}
