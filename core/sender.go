package main

import (
	"fmt"
)

func sender(p *Process) {
	fmt.Println("Starting Sender")
	var pt *Packet = p.create("new IP")
	p.outConn.Send(pt)
}
