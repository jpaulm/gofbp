package testrtn

import "github.com/jpaulm/gofbp"

type ConcatStr struct {
	ipt gofbp.InputArrayConn
	opt gofbp.OutputConn
}

func (concatstr *ConcatStr) Setup(p *gofbp.Process) {
	concatstr.ipt = p.OpenInArrayPort("IN")
	concatstr.opt = p.OpenOutPort("OUT")
}

func (concatstr *ConcatStr) Execute(p *gofbp.Process) {

	for i := 0; i < concatstr.ipt.ArrayLength(); i++ {

		for {
			conn := concatstr.ipt.GetArrayItem(i)
			if conn == nil {
				continue
			}
			var pkt = p.Receive(conn)
			if pkt == nil {
				break
			}
			//fmt.Println("Output: ", pkt.Contents)
			p.Send(concatstr.opt, pkt)
		}
	}

}
