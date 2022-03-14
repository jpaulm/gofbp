package io

import (
	"log"
	"net"

	"github.com/jpaulm/gofbp/core"
)

// https://zetcode.com/golang/socket/

// ReadFromSocket type defines iptIP, ipt, and opt
type ReadFromSocket struct {
	iptIP core.InputConn
	opt   core.OutputConn
}

//Setup method initializes Process
func (readFromSocket *ReadFromSocket) Setup(p *core.Process) {
	readFromSocket.iptIP = p.OpenInPort("PORT")
	readFromSocket.opt = p.OpenOutPort("OUT")
}

//MustRun method
func (ReadFromSocket) MustRun() {}

//Execute method starts Process
func (readFromSocket *ReadFromSocket) Execute(p *core.Process) {

	icpkt := p.Receive(readFromSocket.iptIP)
	port, ok := icpkt.Contents.(string)
	if !ok {
		panic("Parameter (port) not a string")
	}
	p.Discard(icpkt)
	p.Close(readFromSocket.iptIP)

	con, err := net.Dial("tcp", port)

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	for {

		data := make([]byte, 1024)
		_, err = con.Read(data)

		if err != nil {
			return
		}
		pkt := p.Create(string(data))
		p.Send(readFromSocket.opt, pkt)
	}
}
