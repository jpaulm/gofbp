package core

import (
	"fmt"
	"strings"
)

type SubIn struct {
	//ipt   InputConn
	iptIp InputConn
	out   OutputConn
	eipt  InputConn
}

type SubInSS struct {
	SubIn
}

func (subIn *SubIn) Setup(p *Process) {

	//SubIn.ipt = p.OpenInPort("IN")

	subIn.out = p.OpenOutPort("OUT")

	subIn.iptIp = p.OpenInPort("NAME")
}

func (subIn *SubIn) Execute(p *Process) {

	icpkt := p.Receive(subIn.iptIp)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)

	mother := p.network.mother
	subIn.eipt = mother.OpenInPort(param)
	if !strings.Contains(subIn.eipt.(*InPort).portName, ".") {
		subIn.eipt.(*InPort).portName = mother.Name + ":" + subIn.eipt.(*InPort).portName
	}

	for {
		//var pkt = mother.Receive(subIn.eipt)
		var pkt = subIn.eipt.receive(p)
		if pkt == nil {
			trace(p, " Received end of stream")
			break
		}

		pkt.owner = p
		fmt.Println(pkt.Contents)
		p.Send(subIn.out, pkt)

	}

}

func (subIn *SubInSS) Execute(p *Process) {

	icpkt := p.Receive(subIn.iptIp)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)
	mother := p.network.mother
	subIn.eipt = mother.OpenInPort(param)
	if !strings.Contains(subIn.eipt.(*InPort).portName, ":") {
		subIn.eipt.(*InPort).portName = mother.Name + ":" + subIn.eipt.(*InPort).portName
	}

	for {
		//var pkt = mother.Receive(subIn.eipt)
		var pkt = subIn.eipt.receive(p)
		if pkt == nil {
			trace(p, " Received end of stream")
			break
		}

		if pkt.PktType == OpenBracket || pkt.PktType == CloseBracket {
			p.Discard(pkt)
		} else {
			pkt.owner = p
			fmt.Println(pkt.Contents)
			p.Send(subIn.out, pkt)
		}
	}

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

	icpkt := p.Receive(subOut.iptIp)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)
	mother := p.network.mother
	subOut.eopt = mother.OpenOutPort(param)
	if !strings.Contains(subOut.eopt.(*OutPort).portName, ":") {
		subOut.eopt.(*OutPort).portName = mother.Name + ":" + subOut.eopt.(*OutPort).portName
	}

	for {

		pkt := subOut.ipt.receive(p)
		if pkt == nil {
			trace(p, " Received end of stream")
			break
		}
		pkt.owner = p

		fmt.Println(pkt.Contents)

		subOut.eopt.send(p, pkt)

	}

}
