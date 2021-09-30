package main

import (
	"fmt"
	"testing"

	"github.com/jpaulm/gofbp/core"
)

type Bootstrap struct {
	in  core.InputConn
	out core.OutputConn
}

type Forward struct {
	in  core.InputConn
	out core.OutputConn
}

func (c *Bootstrap) Setup(p *core.Process) {
	c.in = p.OpenInPort("IN")
	c.out = p.OpenOutPort("OUT")
}

func (c *Bootstrap) MustRun() {}

func (c *Bootstrap) Execute(p *core.Process) {
	pkt := p.Create("token")
	count := 0
	for {
		p.Send(c.out.(*core.OutPort), pkt)
		pkt = p.Receive(c.in)
		count++
		if count%1000 == 0 {
			fmt.Println(count)
		}
	}
}

func (c *Forward) Setup(p *core.Process) {
	c.in = p.OpenInPort("IN")
	c.out = p.OpenOutPort("OUT")
}

func (c *Forward) Execute(p *core.Process) {
	for {
		pkt := p.Receive(c.in)
		p.Send(c.out.(*core.OutPort), pkt)
	}
}

func TestDetectionBug(t *testing.T) {
	var net *core.Network = core.NewNetwork("DetectionBug")

	boot := net.NewProc("Boot", &Bootstrap{})
	alpha := net.NewProc("Alpha", &Forward{})
	beta := net.NewProc("Beta", &Forward{})

	net.Connect(boot, "OUT", alpha, "IN", 1)
	net.Connect(alpha, "OUT", beta, "IN", 1)
	net.Connect(beta, "OUT", boot, "IN", 1)

	net.Run()
}
