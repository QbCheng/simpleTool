package tcpServer

import (
	"encoding/binary"
)

// CSPacket 客户端和服务器包
type CSPacket struct {
	Header CSPacketHeader
	Body   []byte
}

func NewCSPacket(Version, PassCode uint16, Seq, AppVersion, Cmd uint32, Uid uint64, body []byte) CSPacket {
	ret := CSPacket{
		Header: CSPacketHeader{
			Version:    Version,
			PassCode:   PassCode,
			Seq:        Seq,
			AppVersion: AppVersion,
			Cmd:        Cmd,
			Uid:        Uid,
			BodyLen:    uint32(len(body)),
		},
		Body: body,
	}
	return ret
}

func NewCSPacketByBytes(data []byte) CSPacket {
	head := CSPacketHeader{}
	head.From(data[:head.PacketHeaderLength()])
	ret := CSPacket{
		Header: head,
		Body:   data[head.PacketHeaderLength():],
	}
	return ret
}

func (cs CSPacket) PacketHeaderLength() int {
	return cs.Header.PacketHeaderLength()
}

func (cs CSPacket) PacketBodyLength(header []byte) int {
	return int(binary.BigEndian.Uint32(header[cs.PacketHeaderLength()-4:]))
}

func (cs CSPacket) ToBytes() []byte {
	packet := cs.Header.ToBytes()
	packet = append(packet, cs.Body...)
	return packet
}

// CSPacketHeader 包头
type CSPacketHeader struct {
	Version  uint16
	PassCode uint16
	Seq      uint32

	Uid uint64

	AppVersion uint32
	Cmd        uint32

	BodyLen uint32
}

func (h CSPacketHeader) PacketHeaderLength() int {
	return 28
}

func (h *CSPacketHeader) From(b []byte) {
	pos := 0
	h.Version = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.PassCode = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.Seq = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Uid = binary.BigEndian.Uint64(b[pos:])
	pos += 8
	h.AppVersion = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Cmd = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.BodyLen = binary.BigEndian.Uint32(b[pos:])
	pos += 4
}

func (h *CSPacketHeader) To(b []byte) {
	pos := uintptr(0)
	binary.BigEndian.PutUint16(b[pos:], h.Version)
	pos += 2
	binary.BigEndian.PutUint16(b[pos:], h.PassCode)
	pos += 2
	binary.BigEndian.PutUint32(b[pos:], h.Seq)
	pos += 4
	binary.BigEndian.PutUint64(b[pos:], h.Uid)
	pos += 8
	binary.BigEndian.PutUint32(b[pos:], h.AppVersion)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.Cmd)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.BodyLen)
	pos += 4
}

func (h *CSPacketHeader) ToBytes() []byte {
	bytes := make([]byte, h.PacketHeaderLength())
	h.To(bytes)
	return bytes
}
