// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go

package websocket

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpaulm/gofbp/core"
)

var data_map map[*websocket.Conn][]*core.Packet

type WSRequest struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (wsrequest *WSRequest) Setup(p *core.Process) {
	wsrequest.ipt = p.OpenInPort("ADDR")
	wsrequest.opt = p.OpenOutPort("OUT")
}

var addr *string

var upgrader = websocket.Upgrader{}
var proc *core.Process

//var conn websocket.Conn

func (wsrequest *WSRequest) Execute(p *core.Process) {
	icpkt := p.Receive(wsrequest.ipt)
	path := icpkt.Contents.(string)
	p.Discard(icpkt)
	p.Close(wsrequest.ipt)
	addr = flag.String("addr", path, "http service address")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	proc = p
	err := http.ListenAndServe(*addr, nil)
	log.Fatal(err)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upgrade")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//conn = *c
	defer c.Close()

	if data_map == nil {
		data_map = make(map[*websocket.Conn][]*core.Packet)
	}

	var pkt_list []*core.Packet

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		x := string(message)

		opt := proc.OpenOutPort("OUT")

		//c = &conn
		if x == "@{" {

			pkt_list = make([]*core.Packet, 0)
			data_map[c] = pkt_list
			continue
		}

		if x == "@}" {

			pkt := proc.CreateBracket(core.OpenBracket, "")
			proc.Send(opt, pkt)

			pkt = proc.Create(c)
			proc.Send(opt, pkt)
			for _, pkt := range pkt_list {
				proc.Send(opt, pkt)
			}

			data_map[c] = nil
			pkt = proc.CreateBracket(core.CloseBracket, "")
			proc.Send(opt, pkt)
			continue
		}

		if x == "@kill" {
			time.Sleep(10 * time.Second)
			c.Close()
			break
		}

		pkt := proc.Create(x)
		pkt_list = append(pkt_list, pkt)
		data_map[c] = pkt_list
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
