package io

import (
	"fmt"
	"os"

	"github.com/jpaulm/gofbp/core"
)

type WriteFile struct {
	iptIp *core.InitializationConnection
	ipt   *core.InPort
	opt   *core.OutPort
}

func (writeFile *WriteFile) Setup(p *core.Process) {
	writeFile.iptIp = p.OpenInitializationPort("FILENAME")
	writeFile.ipt = p.OpenInPort("IN")
	writeFile.opt = p.OpenOutPort("OUT", "opt")
}

func (WriteFile) MustRun() {}

func (writeFile *WriteFile) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")
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

		if !writeFile.opt.IsConnected() {
			p.Discard(pkt)
		} else {
			if writeFile.opt == nil {
				panic("WwriteFile - port not specified, but not optional")
			}
			//}
		}

		fmt.Println(p.GetName()+": File", fname, "written")
		//fmt.Println(p.GetName() + " ended")
	}
}
