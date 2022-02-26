package tcpServer

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

type tcpPacketClient struct {
	conn         net.Conn
	id           uint64
	limit, count int
}

func NewTcpPacketClient(id uint64, addr string, limit int) (*tcpPacketClient, error) {
	ret := &tcpPacketClient{
		id:    id,
		limit: limit,
	}
	var err error
	ret.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (cli *tcpPacketClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			packet := NewCSPacket(1, 1, 1, 1, 10001, cli.id, []byte("hello world"))
			_, err := cli.conn.Write(packet.ToBytes())
			if err != nil {
				return
			}
			cli.count++
			if cli.limit > 0 && cli.count >= cli.limit {
				return
			}
		}
	}
}

func TestTcpPacketClient(t *testing.T) {
	client, err := NewTcpPacketClient(1, "127.0.0.1:8080", 100)
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go client.run(&wg)
	wg.Wait()
}

func TestTcpPacketClientMulti(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := uint64(1); i <= 1000; i++ {
		client, err := NewTcpPacketClient(i, "127.0.0.1:8080", 1000)
		if err != nil {
			panic(err)
		}
		go client.run(&wg)
	}
	wg.Wait()
}

func TestTcpPacketClientMultiNotLimit(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := uint64(1); i <= 1000; i++ {
		client, err := NewTcpPacketClient(i, "127.0.0.1:8080", 0)
		if err != nil {
			panic(err)
		}
		go client.run(&wg)
	}
	wg.Wait()
}

type handler struct{}

func (handler) OnConn(conn net.Conn) {
	fmt.Printf("OnConn, {local : %v}, {remote : %v} \n", conn.LocalAddr(), conn.RemoteAddr())
}

func (handler) OnPacket(conn net.Conn, data []byte) {
	packet := NewCSPacketByBytes(data)
	fmt.Printf("OnPacket, {local : %v}, {remote : %v} {header : %#v} {body : %s} \n",
		conn.LocalAddr(), conn.RemoteAddr(), packet.Header, packet.Body)
}

func (handler) OnClose(conn net.Conn) {
	fmt.Printf("OnClose, {local : %v}, {remote : %v} \n", conn.LocalAddr(), conn.RemoteAddr())
}

func TestTcpPacketServer(t *testing.T) {
	server := NewTcpPacketSvr("127.0.0.1", 8080, CSPacket{}, handler{})
	err := server.Server()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Second)
	err = server.Close()
	if err != nil {
		t.Fatal(err)
	}
}
