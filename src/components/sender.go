package components

import (
	core "./core"
	"fmt"
)

func sender(p *core.Process) {
	fmt.Println("Starting Sender")
	var pkt *core.Packet = p.core.create("new IP")
	core.p.outConn.Send(p, pkt)
}
