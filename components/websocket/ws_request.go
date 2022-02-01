// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go

// try https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve

package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpaulm/gofbp/core"
)

//var data_map map[*websocket.Conn][]*core.Packet

type WSRequest struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (wsrequest *WSRequest) Setup(p *core.Process) {
	wsrequest.ipt = p.OpenInPort("ADDR")
	wsrequest.opt = p.OpenOutPort("OUT")
}

var upgrader = websocket.Upgrader{}
var proc *core.Process
var closed_down bool

func (wsrequest *WSRequest) Execute(p *core.Process) {
	icpkt := p.Receive(wsrequest.ipt)
	path := icpkt.Contents.(string)
	p.Discard(icpkt)
	p.Close(wsrequest.ipt)
	proc = p

	log.SetFlags(0)

	log.Printf("main: starting HTTP server")

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	srv := startHttpServer(httpServerExitDone, path)

	for !closed_down {
		//log.Printf("serving for .5 seconds")
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("stopping HTTP server")

	// now close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).  ????????
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	// wait for goroutine started in startHttpServer() to stop
	httpServerExitDone.Wait()

	log.Printf("done. exiting")

}

func startHttpServer(wg *sync.WaitGroup, path string) *http.Server {
	srv := &http.Server{Addr: path}

	http.HandleFunc("/ws", serveWs)

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
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

	var pkt_list []*core.Packet
	var pkt *core.Packet

	opt := proc.OpenOutPort("OUT")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		x := string(message)

		if x == "@{" {

			pkt_list = make([]*core.Packet, 0)

			continue
		}

		if x == "@}" {

			// send out "connection" IPs, then IPs stored in pkt_list ... all surrounded by bracket IPs

			pkt = proc.CreateBracket(core.OpenBracket, "")
			proc.Send(opt, pkt)

			pkt = proc.Create(c) // connection
			proc.Send(opt, pkt)

			for _, pkt := range pkt_list {
				proc.Send(opt, pkt)
			}

			//data_map[c] = nil
			pkt = proc.CreateBracket(core.CloseBracket, "")
			proc.Send(opt, pkt)
			continue
		}

		if x == "@kill" {

			c.Close()
			//if err := srv.Shutdown(context.TODO()); err != nil {
			//	panic(err) // failure/timeout shutting down the server gracefully
			//}
			//break
			closed_down = true
			continue
		}

		pkt = proc.Create(x)
		pkt_list = append(pkt_list, pkt)

	}
}
