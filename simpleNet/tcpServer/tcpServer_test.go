package tcpServer

//
//import (
//	"fmt"
//	"net"
//	"github.com/QbCheng/simpleTool/simpleLogger"
//	"sync"
//	"testing"
//	"time"
//)
//
//type testServer struct {
//	TcpSvr
//	logger simpleLogger.Logger
//}
//
//func NewTestServer(ip string, port int) *testServer {
//	ret := &testServer{}
//	ret.TcpSvr = NewTcpSvr(ip, port, ret)
//	ret.logger = simpleLogger.DefaultLogger()
//
//	err := ret.TcpSvr.Server()
//	if err != nil {
//		ret.logger.Logf("{err : %v}", err)
//	}
//	return ret
//}
//
//func (ts *testServer) OnConn(conn net.Conn) {
//	ts.logger.Logf("OnConn. {local : %v} {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
//}
//
//func (ts *testServer) OnRead(conn net.Conn, data []byte) int {
//	ts.logger.Logf("OnRead. {local : %v} {remote : %v} {msg : %s}", conn.LocalAddr(), conn.RemoteAddr(), data)
//	return len(data)
//}
//
//func (ts *testServer) OnClose(conn net.Conn) {
//	ts.logger.Logf("OnClose. {local : %v} {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
//}
//
//func TestNewTcpSvr( t *testing.T ) {
//	server := NewTestServer("127.0.0.1", 8080)
//	time.Sleep(50 * time.Second)
//	err := server.Close()
//	if err != nil {
//		server.logger.Logf("{err : %v}", err)
//	}
//}
//
//func createClient(limitCount int, wg *sync.WaitGroup) {
//	defer wg.Done()
//	conn, err := net.Dial("tcp", "127.0.0.1:8080")
//	if err != nil {
//		fmt.Println("err : ", err)
//		return
//	}
//	count := 0
//	tick := time.NewTicker(time.Second * 2)
//	for {
//		select {
//		case <-tick.C:
//			count++
//			if _, err := conn.Write([]byte(fmt.Sprintf("hello world %d", count))); err != nil {
//				fmt.Printf("err : %v, local : %v, remote : %v \n", err, conn.LocalAddr(), conn.RemoteAddr())
//				return
//			}
//			if count >= limitCount {
//				if err := conn.Close(); err != nil {
//					fmt.Printf("err : %v, local : %v, remote : %v \n", err, conn.LocalAddr(), conn.RemoteAddr())
//				}
//				return
//			}
//		}
//	}
//}
//
//func TestClient10(t *testing.T) {
//	wg := sync.WaitGroup{}
//	wg.Add(10)
//	for i := 1; i <= 10; i++ {
//		go createClient(10, &wg)
//	}
//	wg.Wait()
//}
//
//
//func createClient2(wg *sync.WaitGroup) {
//	defer wg.Done()
//	conn, err := net.Dial("tcp", "127.0.0.1:8080")
//	if err != nil {
//		fmt.Println("err : ", err)
//		return
//	}
//	count := 0
//	tick := time.NewTicker(time.Millisecond * 200)
//	for {
//		select {
//		case <-tick.C:
//			count++
//			if _, err := conn.Write([]byte(fmt.Sprintf("hello world %d", count))); err != nil {
//				fmt.Printf("err : %v, local : %v, remote : %v \n", err, conn.LocalAddr(), conn.RemoteAddr())
//				return
//			}
//		}
//	}
//}
//
//
//func TestClient100(t *testing.T) {
//	wg := sync.WaitGroup{}
//	wg.Add(100)
//	for i := 1; i <= 100; i++ {
//		go createClient2(&wg)
//	}
//	wg.Wait()
//}
//
//
//func createClientTimeout(wg *sync.WaitGroup) {
//	defer wg.Done()
//	_, err := net.Dial("tcp", "127.0.0.1:8080")
//	if err != nil {
//		fmt.Println("err : ", err)
//		return
//	}
//	time.Sleep(30 * time.Second)
//
//}
//
//func TestClientTimeout(t *testing.T) {
//	wg := sync.WaitGroup{}
//	wg.Add(100)
//	for i := 1; i <= 100; i++ {
//		go createClientTimeout(&wg)
//	}
//	wg.Wait()
//}
