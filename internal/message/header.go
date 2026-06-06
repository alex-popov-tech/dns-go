package message

import (
	"encoding/binary"
	"fmt"
)

func parseHeader(datagram []byte) Header {
	return Header{
		ID:      binary.BigEndian.Uint16(datagram[0:2]),
		Flags:   binary.BigEndian.Uint16(datagram[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(datagram[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(datagram[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(datagram[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(datagram[10:12]),
	}
}

type Header struct {
	// 6 * 2 bytes = 12 bytes
	ID      uint16
	Flags   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h Header) Bytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:], h.ID)
	binary.BigEndian.PutUint16(buf[2:], h.Flags)
	binary.BigEndian.PutUint16(buf[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(buf[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(buf[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(buf[10:], h.ARCOUNT)
	return buf
}

func (h Header) String() string {
	return fmt.Sprintf(
		"ID: %d Flags: %d QDCOUNT: %d ANCOUNT: %d NSCOUNT: %d ARCOUNT: %d",
		h.ID,
		h.Flags,
		h.QDCOUNT,
		h.ANCOUNT,
		h.NSCOUNT,
		h.ARCOUNT,
	)
}

func (h *Header) SetFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint16) {
	h.Flags = qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode
}

func (h Header) Opcode() uint16 {
	return (h.Flags >> 11) & 0xF
}

func (h Header) Rd() uint16 {
	return (h.Flags >> 8) & 0x1
}

func (h Header) Rcode() uint16 {
	return h.Flags & 0xF
}
