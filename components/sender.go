package components

//https://www.geeksforgeeks.org/function-as-a-field-in-golang-structure/

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Sender struct {
	Name    string
	Execute func(*core.Process)
}

func Execute(p *core.Process) {
	fmt.Println("Starting Sender")
	var pkt *core.Packet = p.Create("new IP")
	p.Send(p.OutConn, pkt)
}
