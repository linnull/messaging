package server

import (
	"bufio"
	"lib/packet"
	"log"
	"net"
)

type Proxy struct {
	uuid      uint64
	conn      net.Conn
	rspChan   chan []byte
	closeChan chan bool
	isClosed  bool
}

func NewProxy(uuid uint64, conn net.Conn) *Proxy {
	proxy := new(Proxy)
	proxy.uuid = uuid
	proxy.conn = conn
	proxy.isClosed = false
	proxy.rspChan = make(chan []byte, 1024)
	proxy.closeChan = make(chan bool)
	return proxy
}

func (p *Proxy) recv() {
	reader := bufio.NewReader(p.conn)
	for p.isClosed == false {
		pkg, err := packet.GetPacket(reader)
		if err != nil {
			log.Printf("io.ReadFull err %v\n", err)
			p.Close()
			return
		}
		msg, err := packet.GetRequestMessageFromPacket(pkg)
		if err != nil {
			log.Printf("GetRequestMessageFromPacket err %v\n", err)
			continue
		}
		p.processMessage(msg)
	}
}

func (p *Proxy) send() {
	for p.isClosed == false {
		select {
		case rsp := <-p.rspChan:
			_, err := p.conn.Write(rsp)
			if err != nil {
				log.Printf("send write err %v\n", err)
			}
		case <-p.closeChan:
			return
		}
	}
}

func (p *Proxy) Close() {
	p.isClosed = true
	p.closeChan <- true
	p.conn.Close()
}

func (p *Proxy) processMessage(msg *packet.RequestMessage) {
	log.Println(msg)
}
