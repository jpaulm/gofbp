package components

import (
	"fmt"

	core "github.com/jpaulm/gofbp/src/core"
)

func Sender(p *core.Process) {
	fmt.Println("Starting Sender")
	var pkt *core.Packet = p.Create("new IP")
	p.Send(p.OutConn, pkt)
}
