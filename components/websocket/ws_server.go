//go:build ignore
// +build ignore

package websock

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jpaulm/gofbp/core"
)

type WSServer struct {
	opt core.OutputConn
}

func (wsserver *WSServer) Setup(p *core.Process) {
	wsserver.opt = p.OpenOutPort("OUT")
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
func (wsserver *WSServer) Execute(p *core.Process) {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}
