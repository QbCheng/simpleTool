package tcpserver

import (
	"net"
)

type Packet interface {
	PacketHeaderLength() int
	PacketBodyLength([]byte) int
}

type PacketServerEventCallback interface {
	OnConn(net.Conn)
	OnPacket(net.Conn, []byte)
	OnClose(net.Conn)
}

type TcpPacketSvr struct {
	*TcpSvr

	packet Packet
	cb     PacketServerEventCallback
}

func NewTcpPacketSvr(ip string, port int, packet Packet, cb PacketServerEventCallback) *TcpPacketSvr {
	ret := &TcpPacketSvr{
		packet: packet,
		cb:     cb,
	}
	ret.TcpSvr = NewTcpSvr(ip, port, ret)
	return ret
}

func (s *TcpPacketSvr) OnConn(conn net.Conn) {
	s.cb.OnConn(conn)
}

func (s *TcpPacketSvr) OnRead(conn net.Conn, data []byte) int {
	dataLen := len(data)
	headerLen := s.packet.PacketHeaderLength()
	consumed := 0
	for { // There likely be more than one packet
		if dataLen >= consumed+headerLen { // header is ready
			bodyLen := s.packet.PacketBodyLength(data[consumed : consumed+headerLen])
			if dataLen >= consumed+headerLen+bodyLen { // header and body is ready
				s.cb.OnPacket(conn, data[consumed:consumed+headerLen+bodyLen])
				consumed += headerLen + bodyLen
			} else {
				return consumed
			}
		} else {
			return consumed
		}
	}
}

func (s *TcpPacketSvr) OnClose(conn net.Conn) {
	s.cb.OnClose(conn)
}
