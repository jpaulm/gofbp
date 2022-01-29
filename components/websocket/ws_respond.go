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

func (WSRespond) MustRun() {}

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
	conn, _ := pkt.Contents.(*websocket.Conn)

	p.Discard(pkt)

	for {
		pkt = p.Receive(wsrespond.ipt)
		if pkt.PktType == core.CloseBracket {
			break
		}
		data, ok := pkt.Contents.(string)
		if !ok {
			log.Println("write: data not string")
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			p.Discard(pkt)
		}
	}
	p.Discard(pkt)
}
