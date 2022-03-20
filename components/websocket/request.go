package websocket

// The problem was to pass an extra parameter to a Handler function - I found several packages described on StackOverflow
// which were designed for this task, but I found them all confusing!  Then I found Alex's answer, which
// was simple and easy to use, needed no additional packages, and, above all, was CLEAR!

// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
// https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpaulm/gofbp/core"
)

type WSRequest struct {
	ipt         core.InputConn
	opt         core.OutputConn
	proc        *core.Process
	closed_down int32
}

func (wsrequest *WSRequest) Setup(p *core.Process) {
	wsrequest.ipt = p.OpenInPort("ADDR")
	wsrequest.opt = p.OpenOutPort("OUT")
}

var upgrader = websocket.Upgrader{}

// var proc *core.Process
// var closed_down bool

func (wsrequest *WSRequest) Execute(p *core.Process) {
	icpkt := p.Receive(wsrequest.ipt)
	path := icpkt.Contents.(string)
	p.Discard(icpkt)
	p.Close(wsrequest.ipt)
	wsrequest.proc = p

	log.SetFlags(0)

	log.Printf("main: starting HTTP server")

	wg := &sync.WaitGroup{}

	wg.Add(1)

	//srv := startHttpServer(wg, path, wsrequest)
	srv := &http.Server{Addr: path}
	mh := myHandler{wsr: wsrequest}

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := http.ListenAndServe(path, mh); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	//return srv

	for {
		if atomic.CompareAndSwapInt32(&wsrequest.closed_down, 1, 1) {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("stopping HTTP server")

	// now close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).  ????????
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	wg.Done()

	// wait for goroutine started in startHttpServer() to stop
	wg.Wait()

	log.Printf("done. exiting")

}

type myHandler struct {
	wsr *WSRequest
}

// ServeHTTP ensures that myHandler is an instance of Handler interface{}

func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upgrade")
	//fmt.Printf("%+v\n", w)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	var pkt_list []*core.Packet
	var pkt *core.Packet

	//var cn websocket.Conn
	cn := *c
	c2 := &cn
	opt := mh.wsr.proc.OpenOutPort("OUT")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		x := string(message)

		if x == "@kill" {
			pkt = mh.wsr.proc.CreateSignal(x)
			mh.wsr.proc.Send(opt, pkt)
			atomic.StoreInt32(&mh.wsr.closed_down, 1)
			break
		}

		if x == "@{" {

			pkt_list = make([]*core.Packet, 0)

			continue
		}

		if x == "@}" {

			// send out "connection" IPs, then IPs retrieved from pkt_list ... all surrounded by bracket IPs

			pkt = mh.wsr.proc.CreateBracket(core.OpenBracket, "")
			mh.wsr.proc.Send(opt, pkt)

			pkt = mh.wsr.proc.Create(c2) // connection
			mh.wsr.proc.Send(opt, pkt)

			for _, pkt := range pkt_list {
				mh.wsr.proc.Send(opt, pkt)
			}

			//data_map[c] = nil
			pkt = mh.wsr.proc.CreateBracket(core.CloseBracket, "")
			mh.wsr.proc.Send(opt, pkt)
			continue
		}

		pkt = mh.wsr.proc.Create(x)
		pkt_list = append(pkt_list, pkt)
	}
}
