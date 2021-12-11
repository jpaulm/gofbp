/*Package io implements gofbp I/O*/
package io

import (
	//"core"
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

	icpkt := p.Receive(readFile.ipt)
	fname, ok := icpkt.Contents.(string)
	if !ok {
		panic("Parameter (file name) not a string")
	}
	f, err := os.Open(fname)
	if err != nil {
		panic("Unable to read file: " + fname)
	}
	p.Discard(icpkt)
	p.Close(readFile.ipt)

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
		p.Send(readFile.opt, pkt)
		rec = rec[i+1:]
	}

}
