package components

import (
	"fmt"
	"github.com/jpaulm/gofbp/core"
)

var Name string = "Sender"

func Execute(p *core.Process) {
	fmt.Println("Sender started")
	var pkt *core.Packet = p.Create("new IP")
	p.Send(p.OutConn, pkt)
	pkt = p.Create("2nd IP")
	p.Send(p.OutConn, pkt)
	fmt.Println("Sender ended")
}
