package main

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

type Forward struct {
	Limit int

	in  core.InputConn
	out core.OutputConn
}

func (c *Forward) Setup(p *core.Process) {
	c.in = p.OpenInPort("IN")
	c.out = p.OpenOutPort("OUT")
}

func (c *Forward) Execute(p *core.Process) {
	limit := 0
	for {
		pkt := p.Receive(c.in)
		if pkt == nil {
			return
		}
		if limit++; c.Limit > 0 && limit > c.Limit {
			p.Discard(pkt)
			return
		}
		p.Send(c.out, pkt)
	}
}

func TestForwarding(t *testing.T) {
	net := core.NewNetwork("Forwarding")

	kick := net.NewProc("Kick", &testrtn.Kick{})
	alpha := net.NewProc("Alpha", &Forward{Limit: 5})
	beta := net.NewProc("Beta", &Forward{})
	gamma := net.NewProc("Gamma", &Forward{})

	net.Connect(kick, "OUT", alpha, "IN", 1)
	net.Connect(alpha, "OUT", beta, "IN", 1)
	net.Connect(beta, "OUT", gamma, "IN", 1)
	net.Connect(gamma, "OUT", alpha, "IN", 1)

	net.Run()
}
