package components

import (
	"fmt"

	"github.com/jpaulm/gofbp/src/core"
)

func sender(p *core.Process) {
	fmt.Println("Starting Sender")
	var pkt *core.Packet = p.Create("new IP")
	p.Send(p.OutConn, pkt)
}
