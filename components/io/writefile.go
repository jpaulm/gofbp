package io

import (
	"fmt"
	"os"

	"github.com/jpaulm/gofbp/core"
)

type WriteFile struct {
	iptIP core.InputConn
	ipt   core.InputConn
	opt   core.OutputConn
}

func (writeFile *WriteFile) Setup(p *core.Process) {
	writeFile.iptIP = p.OpenInPort("FILENAME")
	writeFile.ipt = p.OpenInPort("IN")
	writeFile.opt = p.OpenOutPortOptional("OUT")
}

func (WriteFile) MustRun() {}

func (writeFile *WriteFile) Execute(p *core.Process) {

	icpkt := p.Receive(writeFile.iptIP)
	fname, ok := icpkt.Contents.(string)
	if !ok {
		panic("Parameter (file name) not a string")
	}
	p.Discard(icpkt)
	p.Close(writeFile.iptIP)

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
			p.Send(writeFile.opt, pkt)
		}

	}
	fmt.Println(p.Name+": File", fname, "written")
}
