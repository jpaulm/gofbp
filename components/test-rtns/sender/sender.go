package sender

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

var Name string = "Sender"

// var opt *core.OutPort

//type Sender struct{ opt *core.OutPort }

func OpenPorts(p *core.Process) {
	opt = p.OpenOutPort("OUT")
}

func Execute(p *core.Process) {
	fmt.Println(p.Name + " started")
	var pkt *core.Packet
	for i := 0; i < 15; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i) + " (" + p.Name + ")")
		p.Send(opt.Conn, pkt)
	}
	fmt.Println(p.Name + " ended")
}
