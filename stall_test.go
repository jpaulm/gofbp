package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jpaulm/gofbp/core"
)

type Either struct {
	in   core.InputConn
	odd  core.OutputConn
	even core.OutputConn
}

func (c *Either) Setup(p *core.Process) {
	c.in = p.OpenInPort("IN")
	c.odd = p.OpenOutPort("ODD")
	c.even = p.OpenOutPort("EVEN")
}

func (c *Either) Execute(p *core.Process) {
	for {
		pkt := p.Receive(c.in)
		if pkt == nil {
			return
		}

		v, err := strconv.Atoi(pkt.Contents.(string))
		if err != nil {
			p.Discard(pkt)
			continue
		}

		if v%2 == 0 {
			p.Send(c.even, pkt)
		} else {
			p.Send(c.odd, pkt)
		}
	}
}

type Printer struct {
	prefix string
	in     core.InputConn
}

func (c *Printer) Setup(p *core.Process) {
	c.in = p.OpenInPort("IN")
}

func (c *Printer) Execute(p *core.Process) {
	for {
		pkt := p.Receive(c.in)
		if pkt == nil {
			return
		}
		fmt.Println(c.prefix, ":", pkt.Contents)
		p.Discard(pkt)
	}
}

func TestStallBug(t *testing.T) {
	net := core.NewNetwork("StallBug")

	either := net.NewProc("Either", &Either{})
	odd := net.NewProc("Odd", &Printer{prefix: "ODD"})
	even := net.NewProc("Even", &Printer{prefix: "EVEN"})

	net.Initialize("0", either, "IN")
	net.Connect(either, "ODD", odd, "IN", 1)
	net.Connect(either, "EVEN", even, "IN", 1)

	net.Run()
}
