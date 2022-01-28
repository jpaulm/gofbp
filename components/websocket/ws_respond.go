// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go

package websocket

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/jpaulm/gofbp/core"
)

type WSRespond struct {
	ipt core.InputConn
}

func (wsrespond *WSRespond) Setup(p *core.Process) {
	wsrespond.ipt = p.OpenInPort("IN")
}

func (wsrespond *WSRespond) Execute(p *core.Process) {
	pkt := p.Receive(wsrespond.ipt) // open bracket
	if pkt == nil {
		return
	}
	p.Discard(pkt)

	pkt = p.Receive(wsrespond.ipt) // connection
	if pkt == nil {
		return
	}
	conn, ok := pkt.Contents.(*websocket.Conn)
	if !ok {
		i := 1
		i++
	}

	p.Discard(pkt)

	for {
		pkt := p.Receive(wsrespond.ipt)
		if pkt.PktType == core.CloseBracket {
			break
		}
		data, _ := pkt.Contents.([]byte)

		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			p.Discard(pkt)
		}
	}
	p.Discard(pkt)
}
