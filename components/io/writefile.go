package io

import (
	"fmt"
	"os"

	"github.com/jpaulm/gofbp/core"
)

type WriteFile struct {
	iptIp core.InputConn
	ipt   core.InputConn
	opt   core.OutputConn
}

func (writeFile *WriteFile) Setup(p *core.Process) {
	writeFile.iptIp = p.OpenInPort("FILENAME")
	writeFile.ipt = p.OpenInPort("IN")
	writeFile.opt = p.OpenOutPort("OUT", "opt")
}

func (WriteFile) MustRun() {}

func (writeFile *WriteFile) Execute(p *core.Process) {
	fmt.Println(p.GetName() + " started")
	icpkt := p.Receive(writeFile.iptIp)
	fname := icpkt.Contents.(string)
	p.Discard(icpkt)

	f, err := os.Create(fname)
	if err != nil {
		panic("Unable to open file: " + fname)
	}
	defer f.Close()

	for {
		var pkt = p.Receive(writeFile.ipt)
		if pkt == nil {
			break
		}

		data := []byte(pkt.Contents.(string) + "\n")

		_, err2 := f.Write(data)

		if err2 != nil {
			panic("Unable to write file: " + fname)
		}

		if writeFile.opt.GetType() == "OutPort" {
			p.Send(writeFile.opt.(*core.OutPort), pkt)
		} else {
			p.Discard(pkt)
		}
	}

	fmt.Println(p.GetName()+": File", fname, "written")
	fmt.Println(p.GetName() + " ended")
}
