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
	pkt := p.Receive(wsrespond.ipt)
	for {

		if pkt == nil {
			return
		}

		p.Discard(pkt)                 // discard open bracket
		pkt = p.Receive(wsrespond.ipt) // connection
		if pkt == nil {
			return
		}
		conn, ok := pkt.Contents.(*websocket.Conn)
		if !ok {
			panic("WSRespond - IP after open bracket not *websocket.Conn")
		}
		p.Discard(pkt)
		err := conn.WriteMessage(websocket.TextMessage, []byte("@{"))
		if err != nil {
			log.Println("write:", err)
			break
		}
		pkt = p.Receive(wsrespond.ipt)

		for {
			if pkt.PktType == core.CloseBracket {
				p.Discard(pkt)
				err := conn.WriteMessage(websocket.TextMessage, []byte("@}"))
				if err != nil {
					log.Println("write:", err)
					break
				}
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
		}
		if pkt.PktType == core.Signal && pkt.Contents.(string) == "@kill" {
			conn.Close()
			p.Discard(pkt)
			pkt = p.Receive(wsrespond.ipt)
		}
	}
}
