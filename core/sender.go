package main

import (
	"fmt"
)

func sender(p *Process) {
	fmt.Println("Test")
	var pt Packet = p.create("new IP")
	p.Send(pt)
}
