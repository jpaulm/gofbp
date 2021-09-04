package components

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

var Name string = "Sender"

func Execute(p *core.Process) {
	fmt.Println("Sender started")
	var pkt *core.Packet
	for i := 0; i < 25; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(p.OutConn, pkt)
	}
	fmt.Println("Sender ended")
}
