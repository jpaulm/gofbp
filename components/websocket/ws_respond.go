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
	for {
		pkt := p.Receive(wsrespond.ipt)
		if pkt == nil {
			return
		}
		if pkt.PktType == core.OpenBracket {
			p.Discard(pkt)
			pkt = p.Receive(wsrespond.ipt) // connection
			if pkt == nil {
				return
			}
			conn, _ = pkt.Contents.(websocket.Conn)
			p.Discard(pkt)
			continue
		}

		if pkt.PktType == core.CloseBracket {
			p.Discard(pkt)
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(pkt.Contents.(string)))
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			p.Discard(pkt)
		}
	}
}
