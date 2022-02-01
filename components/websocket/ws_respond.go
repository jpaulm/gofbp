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

var conn websocket.Conn

func (wsrespond *WSRespond) Execute(p *core.Process) {
	pkt := p.Receive(wsrespond.ipt)
	for {

		if pkt == nil {
			return
		}
		if pkt.PktType != core.OpenBracket {
			panic("WSRespond - first IP not open bracket")
		}
		p.Discard(pkt)
		pkt = p.Receive(wsrespond.ipt) // connection
		if pkt == nil {
			return
		}
		conn, _ = pkt.Contents.(websocket.Conn)
		p.Discard(pkt)
		pkt = p.Receive(wsrespond.ipt)

		for {
			if pkt.PktType == core.CloseBracket {
				p.Discard(pkt)
				pkt = p.Receive(wsrespond.ipt)
				break
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(pkt.Contents.(string)))
			if err != nil {
				log.Println("write:", err)
				break
			} else {
				p.Discard(pkt)
			}
			pkt = p.Receive(wsrespond.ipt)

			//pkt = p.Receive(wsrespond.ipt)
		}
		if pkt.Contents.(string) == "@kill" {
			conn.Close()

		}
		pkt = p.Receive(wsrespond.ipt)
	}
}
