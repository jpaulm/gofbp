package io

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/jpaulm/gofbp/core"
)

// https://stackoverflow.com/questions/42126978/how-can-i-keep-reading-using-net-conn-read-method

// ReadFromSocket type defines iptIP, ipt, and opt
type ReadFromSocket struct {
	iptIP core.InputConn
	opt   core.OutputConn
	p     *core.Process
	rc    bool
	wg    *sync.WaitGroup
}

func (readFromSocket *ReadFromSocket) Setup(p *core.Process) {
	readFromSocket.iptIP = p.OpenInPort("PORT")
	readFromSocket.opt = p.OpenOutPort("OUT")
	readFromSocket.p = p
	readFromSocket.wg = &sync.WaitGroup{}
}

func (ReadFromSocket) MustRun() {}

func (readFromSocket *ReadFromSocket) Execute(p *core.Process) {

	icpkt := p.Receive(readFromSocket.iptIP)
	port, ok := icpkt.Contents.(string)
	if !ok {
		panic("Parameter (port) not a string")
	}
	p.Discard(icpkt)
	p.Close(readFromSocket.iptIP)

	//con, err := net.Dial("tcp", port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		readFromSocket.rc = false

		readFromSocket.wg.Add(1)
		go readFromSocket.handleRequest(conn)
		readFromSocket.wg.Wait()
		if readFromSocket.rc {
			break
		}
	}
}

// Handles incoming requests.
func (rfs *ReadFromSocket) handleRequest(conn net.Conn) {
	defer rfs.wg.Done()
	defer conn.Close()
	buffer := make([]byte, 1024)
	//j := 0
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			log.Println(err)
			rfs.rc = err == io.EOF // true if EOF
			return
		}
		//j := 1
		if n > 0 || err == io.EOF {
			//data := string(buffer[:n])
			message := string(buffer[:n])
			pkt := rfs.p.Create(message)
			rfs.p.Send(rfs.opt, pkt)
			_, err = conn.Write([]byte("\u0006"))
			/*
				for {
					for i := 1; i < n; i++ {
						if data[i:i+1] == "}" {
							message := data[j:i]
							pkt := rfs.p.Create(message)
							rfs.p.Send(rfs.opt, pkt)
							i = i + 2
							// j = i
						}
					}
				}
			*/
		}
		//if rfs.rc {
		//	return
		//}
	}

}
