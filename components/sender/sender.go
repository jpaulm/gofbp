package components

import (
	"fmt"
	"github.com/jpaulm/gofbp/core"
	"strconv"
)

var Name string = "Sender"

func Execute(p *core.Process) {
	fmt.Println("Sender started")
	var pkt *core.Packet
	for i := 0; i < 6; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(p.OutConn, pkt)
	}
	fmt.Println("Sender ended")
}
