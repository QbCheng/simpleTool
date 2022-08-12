package tcpserver

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"simpleTool/simpleLogger"
	"strconv"
	"sync"
	"time"
)

// ConnEventCallback 连接事件回调
type ConnEventCallback interface {
	OnConn(net.Conn)
	OnRead(net.Conn, []byte) int
	OnClose(net.Conn)
}

type TcpConnInfo struct {
	chanWrite chan []byte
}

type Option func(*TcpSvr)

const (
	defaultTcpReadTimeout  = 30 * time.Second
	defaultTcpWriteTimeout = 10 * time.Second
	defaultReadBufferSize  = 1024 * 10
)

var (
	ErrorConnectionNotExist           = errors.New(" Tcp server connection doesn't exist. ")
	ErrorWriteConnectionBufferTimeout = errors.New(" Tcp server write connection buffer timeout. ")
)

func WithTcpReadTimeout(tcpReadTimeout time.Duration) Option {
	return func(tcpServer *TcpSvr) {
		tcpServer.tcpReadTimeout = tcpReadTimeout
	}
}

func WithTcpWriteTimeout(tcpWriteTimeout time.Duration) Option {
	return func(tcpServer *TcpSvr) {
		tcpServer.tcpWriteTimeout = tcpWriteTimeout
	}
}

func WithLogger(logger simpleLogger.Logger) Option {
	return func(tcpServer *TcpSvr) {
		tcpServer.logger = logger
	}
}

type TcpSvr struct {
	cb   ConnEventCallback
	ip   string
	port int

	listener       net.Listener
	lockOfConnInfo sync.RWMutex
	mapOfConnInfo  map[net.Conn]TcpConnInfo

	tcpReadTimeout  time.Duration
	tcpWriteTimeout time.Duration
	readBufferSize  int
	logger          simpleLogger.Logger
}

func NewTcpSvr(ip string, port int, cb ConnEventCallback, options ...Option) *TcpSvr {
	ret := &TcpSvr{
		cb:              cb,
		ip:              ip,
		port:            port,
		mapOfConnInfo:   map[net.Conn]TcpConnInfo{},
		tcpReadTimeout:  defaultTcpReadTimeout,
		tcpWriteTimeout: defaultTcpWriteTimeout,
		readBufferSize:  defaultReadBufferSize,
		logger: simpleLogger.NewLogger(
			simpleLogger.WithCallPath(3),
			simpleLogger.WithFlag(log.Lshortfile|log.LstdFlags),
		),
	}
	for i := range options {
		options[i](ret)
	}
	return ret
}

func (s *TcpSvr) err(layout string, v ...interface{}) {
	t := fmt.Sprintf(layout, v...)
	s.logger.Logf("[tcpSvr err] " + t)
}

func (s *TcpSvr) info(layout string, v ...interface{}) {
	t := fmt.Sprintf(layout, v...)
	s.logger.Logf("[tcpSvr err] " + t)
}

func (s *TcpSvr) Server() error {
	var err error
	addr := s.ip + ":" + strconv.Itoa(s.port)
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go s.accept(s.listener)
	return nil
}

func (s *TcpSvr) accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.err("Error accepting. {err : %v}", err)
			return
		}

		s.lockOfConnInfo.Lock()
		if s.mapOfConnInfo == nil {
			s.err("TcpSvr closing. {err : %v}", err)
			return
		}
		s.mapOfConnInfo[conn] = TcpConnInfo{chanWrite: make(chan []byte, 1)}
		s.lockOfConnInfo.Unlock()
		s.info("New Connection. {local : %v} {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
		s.cb.OnConn(conn)
		go s.connRead(conn)
		go s.connWrite(conn, s.mapOfConnInfo[conn].chanWrite)
	}
}

// 读 协程
func (s *TcpSvr) connRead(conn net.Conn) {
	var buff bytes.Buffer
	readBuf := make([]byte, s.readBufferSize)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(s.tcpReadTimeout))
		readLen, err := conn.Read(readBuf)
		if err == nil {
			buff.Write(readBuf[0:readLen])
			consumedLen := s.cb.OnRead(conn, buff.Bytes())
			if consumedLen > 0 {
				buff.Next(consumedLen)
			}
		} else {
			if errors.Is(err, net.ErrClosed) {
				s.info("Connect closed. {local : %v} {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
				break
			}
			if err == io.EOF {
				s.info("conn read {error:io.EOF}. {local : %v} {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
				break
			} else {
				s.err("error occurs when read from tcp.  {local : %v} {remote : %v} {err : %v}", conn.LocalAddr(), conn.RemoteAddr(), err)
				break
			}
		}
	}
	s.cb.OnClose(conn)
	s.destroyConn(conn)
}

// 写 协程
func (s *TcpSvr) connWrite(conn net.Conn, chanWrite <-chan []byte) {
	for {
		select {
		case writeData, ok := <-chanWrite:
			if !ok { // chan is closed
				s.info("chan is closed {local : %v}, {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
				break
			}
			if writeData == nil {
				s.info(" External closed. {local : %v}, {remote : %v}", conn.LocalAddr(), conn.RemoteAddr())
				break
			}

			// 设置 写操作的超时时间
			err := conn.SetWriteDeadline(time.Now().Add(s.tcpWriteTimeout))
			if err != nil {
				s.err("set write deadline error. {err : %v}", err)
				_ = conn.Close()
			}

			sentLen, err := conn.Write(writeData)
			if sentLen < len(writeData) || err != nil {
				//todo: retry?
				s.err("Failed to write tcp data. {local : %v} {remote : %v} {err : %v}, {dateLen : %v}, {writeLen : %v}",
					conn.LocalAddr(),
					conn.RemoteAddr(),
					err,
					len(writeData),
					sentLen,
				)
				_ = conn.Close()
				break
			}
		}
	}
}

func (s *TcpSvr) destroyConn(conn net.Conn) {
	s.lockOfConnInfo.Lock()
	defer s.lockOfConnInfo.Unlock()

	if info, exists := s.mapOfConnInfo[conn]; exists {
		info.chanWrite <- nil
		delete(s.mapOfConnInfo, conn)
	}
}

func (s *TcpSvr) WriteData(conn net.Conn, data []byte) error {
	var chanWrite chan []byte = nil

	s.lockOfConnInfo.RLock()
	info, exists := s.mapOfConnInfo[conn]
	if !exists {
		s.lockOfConnInfo.RUnlock()
		return ErrorConnectionNotExist
	}
	chanWrite = info.chanWrite
	s.lockOfConnInfo.RUnlock()

	t := time.NewTimer(3 * time.Second)
	defer t.Stop()
	select {
	case chanWrite <- data:
	case <-t.C:
		return ErrorWriteConnectionBufferTimeout
	}

	return nil
}

func (s *TcpSvr) Close() error {

	if err := s.listener.Close(); err != nil {
		s.err(" tcp server listener close. {err : %v} ", err)
	}

	s.lockOfConnInfo.RLock()
	for conn := range s.mapOfConnInfo {
		if err := conn.Close(); err != nil {
			s.err(" tcp server connection close.  {local : %v} {remote : %v} {err : %v} ", conn.LocalAddr(), conn.RemoteAddr(), err)
		}
	}
	s.mapOfConnInfo = map[net.Conn]TcpConnInfo{}
	s.lockOfConnInfo.RUnlock()
	return nil
}
