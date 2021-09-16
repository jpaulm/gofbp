package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type ConcatStr struct {
	ipt core.InputConn
	//opt     *core.OutPort
	opt     core.OutputConn
	MustRun bool
}

func (concatstr *ConcatStr) OpenPorts(p *core.Process) {
	concatstr.ipt = p.OpenInPort("IN")
	concatstr.opt = p.OpenOutPort("OUT")
}

func (concatstr *ConcatStr) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")

	for i := 0; i < concatstr.ipt.ArrayLength(); i++ {

		for {
			conn := concatstr.ipt.GetArrayItem(i)
			if conn == nil {
				continue
			}
			var pkt = p.Receive(concatstr.ipt.GetArrayItem(i))
			if pkt == nil {
				break
			}
			//fmt.Println("Output: ", pkt.Contents)
			p.Send(concatstr.opt.Conn, pkt)
		}
	}
	fmt.Println(p.Name + " ended")
}

func (concatstr *ConcatStr) GetMustRun() bool {
	return false
}
