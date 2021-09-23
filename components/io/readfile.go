package io

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jpaulm/gofbp/core"
)

type ReadFile struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (readFile *ReadFile) Setup(p *core.Process) {
	readFile.ipt = p.OpenInPort("FILENAME")
	readFile.opt = p.OpenOutPort("OUT")
}

func (readFile *ReadFile) Execute(p *core.Process) {
	fmt.Println(p.GetName() + " started")
	icpkt := p.Receive(readFile.ipt)
	f, err := os.Open(icpkt.Contents.(string))
	if err != nil {
		panic("Unable to read file: " + f.Name())
	}
	p.Discard(icpkt)

	var pkt *core.Packet
	var rec string

	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			//fmt.Println(string(buf[:n]))
			rec += string(buf[:n])
		}
	}

	for {
		i := strings.Index(rec, "\n")
		if i == -1 {
			break
		}
		pkt = p.Create(rec[:i])
		p.Send(readFile.opt.(*core.OutPort), pkt)
		rec = rec[i+1:]
	}
	fmt.Println(p.GetName() + " ended")
}
