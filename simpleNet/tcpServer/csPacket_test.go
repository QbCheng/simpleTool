package tcpserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeader(t *testing.T) {
	header := CSPacketHeader{
		Version:    1,
		PassCode:   1,
		Seq:        1,
		Uid:        100001,
		AppVersion: 1,
		Cmd:        1,
		BodyLen:    16,
	}
	hData := make([]byte, header.PacketHeaderLength())
	header.To(hData)
	newHeader := CSPacketHeader{}
	newHeader.From(hData)
	assert.Equal(t, newHeader, header)
}
