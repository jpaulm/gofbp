package io

import (
	"fmt"
	"log"
	"net"

	"github.com/jpaulm/gofbp/core"
)

// https://zetcode.com/golang/socket/

// WriteToSocket type defines iptIP, ipt, and opt
type WriteToSocket struct {
	iptIP core.InputConn
	ipt   core.InputConn
	opt   core.OutputConn
}

//Setup method initializes Process
func (writeToSocket *WriteToSocket) Setup(p *core.Process) {
	writeToSocket.iptIP = p.OpenInPort("PORT")
	writeToSocket.ipt = p.OpenInPort("IN")
	writeToSocket.opt = p.OpenOutPortOptional("OUT")
}

//MustRun method
func (WriteToSocket) MustRun() {}

//Execute method starts Process
func (writeToSocket *WriteToSocket) Execute(p *core.Process) {

	icpkt := p.Receive(writeToSocket.iptIP)
	port, ok := icpkt.Contents.(string)
	if !ok {
		panic("Parameter (port) not a string")
	}
	p.Discard(icpkt)
	p.Close(writeToSocket.iptIP)

	con, err := net.Dial("tcp", port)

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	for {
		var pkt = p.Receive(writeToSocket.ipt)
		if pkt == nil {
			break
		}
		data := fmt.Sprint(pkt.Contents)
		_, err = con.Write([]byte(data))

		if err != nil {
			log.Fatal(err)
		}

		if !writeToSocket.opt.IsConnected() {
			p.Discard(pkt)
		} else {
			p.Send(writeToSocket.opt, pkt)
		}
	}
}
