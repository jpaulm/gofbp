package sender

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

var Name string = "Sender"

func Execute(p *core.Process) {
	fmt.Println(p.Name + " started")
	var pkt *core.Packet
	for i := 0; i < 15; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i) + " (" + p.Name + ")")
		p.Send(p.OutConn, pkt)
	}
	fmt.Println(p.Name + " ended")
}
