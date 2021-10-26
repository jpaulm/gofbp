package core

import (
	"fmt"
	//"github.com/jpaulm/gofbp/src/core"
)

type SubIn struct {
	//ipt   InputConn
	iptIp InputConn
	out   OutputConn
	eipt  InputConn
}

func (subIn *SubIn) Setup(p *Process) {

	//SubIn.ipt = p.OpenInPort("IN")

	subIn.out = p.OpenOutPort("OUT")

	subIn.iptIp = p.OpenInPort("NAME")
}

//func (SubIn) MustRun() {}

func (subIn *SubIn) Execute(p *Process) {
	//fmt.Println(p.GetName() + " started")

	icpkt := p.Receive(subIn.iptIp)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)
	mother := p.network.Mother
	subIn.eipt = mother.OpenInPort(param)

	for {
		//var pkt = mother.Receive(subIn.eipt)
		var pkt = subIn.eipt.receive(p)
		if pkt == nil {
			break
		}
		pkt.owner = p

		fmt.Println(pkt.Contents)

		p.Send(subIn.out, pkt)

	}

	//fmt.Println(p.GetName() + " ended")
}

type SubOut struct {
	ipt   InputConn
	iptIp InputConn
	eopt  OutputConn
}

func (subOut *SubOut) Setup(p *Process) {

	subOut.ipt = p.OpenInPort("IN")

	subOut.iptIp = p.OpenInPort("NAME")
}

//func (SubOut) MustRun() {}

func (subOut *SubOut) Execute(p *Process) {
	//fmt.Println(p.GetName() + " started")

	icpkt := p.Receive(subOut.iptIp)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)
	mother := p.network.Mother
	subOut.eopt = mother.OpenOutPort(param)

	for {
		//var pkt = mother.Receive(subIn.eipt)
		var pkt = subOut.ipt.receive(p)
		if pkt == nil {
			break
		}
		pkt.owner = p

		fmt.Println(pkt.Contents)
		//pkt.owner = mother

		//p.Send(subOut.out, pkt)
		subOut.eopt.send(p, pkt)

	}

	//fmt.Println(p.GetName() + " ended")
}
