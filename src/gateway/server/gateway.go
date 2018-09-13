package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/micro/go-micro"
)

type Gateway struct {
	service  micro.Service
	listener net.Listener
	port     uint16
	running  int32
	GUID     uint64
}

func NewGateway() *Gateway {
	gw := new(Gateway)
	return gw
}

func (gw *Gateway) Init(service micro.Service) (err error) {
	atomic.StoreInt32(&gw.running, 0)
	gw.service = service

	gw.port = uint16(20000)
	gw.listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", gw.port))
	if err != nil {
		return err
	}
	log.Printf("tcp server listening in %d\n", gw.port)

	return nil
}

func (gw *Gateway) Run() {
	atomic.StoreInt32(&gw.running, 1)
	go gw.runServer()
}

func (gw *Gateway) runServer() {
	for atomic.LoadInt32(&gw.running) == 1 {
		conn, err := gw.listener.Accept()
		if err == nil {
			go gw.ProcessNewConn(conn)
		} else {
			log.Printf("Accept failed: %s; retrying in 1s\n", gw.port)
			time.Sleep(time.Second)
		}

	}
}

func (gw *Gateway) ProcessNewConn(conn net.Conn) {
	if tc, ok := conn.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(10 * time.Second)
	}

	newUUID := atomic.AddUint64(&gw.GUID, 1)

	proxy := NewProxy(newUUID, conn)
	go proxy.recv()
	proxy.send()
}
