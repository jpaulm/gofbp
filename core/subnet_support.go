package core

import (
	//"fmt"
	"strings"
)

type SubIn struct {
	//ipt   InputConn
	iptIP InputConn
	out   OutputConn
	eipt  InputConn
}

type SubInSS struct {
	SubIn
}

func (subIn *SubIn) Setup(p *Process) {

	//SubIn.ipt = p.OpenInPort("IN")

	subIn.out = p.OpenOutPort("OUT")

	subIn.iptIP = p.OpenInPort("NAME")
}

func (subIn *SubIn) Execute(p *Process) {

	icpkt := p.Receive(subIn.iptIP)
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
			//trace(p, " Received end of stream")
			break
		}

		pkt.owner = p
		//fmt.Println(pkt.Contents)
		p.Send(subIn.out, pkt)

	}

}

func (subIn *SubInSS) Execute(p *Process) {

	icpkt := p.Receive(subIn.iptIP)
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
			//trace(p, " Received end of stream")
			break
		}

		if pkt.PktType == OpenBracket {
			p.Discard(pkt)
		} else {
			if pkt.PktType != CloseBracket {
				pkt.owner = p
				//fmt.Println(pkt.Contents)
				p.Send(subIn.out, pkt)
			} else {
				p.Discard(pkt)
				//subIn.out.Close()
				return
			}
		}
	}
}

type SubOut struct {
	ipt   InputConn
	iptIP InputConn
	eopt  OutputConn
}

type SubOutSS struct {
	SubOut
}

func (subOut *SubOut) Setup(p *Process) {

	subOut.ipt = p.OpenInPort("IN")

	subOut.iptIP = p.OpenInPort("NAME")
}

//func (SubOut) MustRun() {}

func (subOut *SubOut) Execute(p *Process) {

	icpkt := p.Receive(subOut.iptIP)
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
			//trace(p, " Received end of stream")
			break
		}
		pkt.owner = p

		//fmt.Println(pkt.Contents)

		subOut.eopt.send(p, pkt)
	}
}
func (subOut *SubOutSS) Execute(p *Process) {

	icpkt := p.Receive(subOut.iptIP)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)
	mother := p.network.mother
	subOut.eopt = mother.OpenOutPort(param)
	if !strings.Contains(subOut.eopt.(*OutPort).portName, ":") {
		subOut.eopt.(*OutPort).portName = mother.Name + ":" + subOut.eopt.(*OutPort).portName
	}

	pkt := p.CreateBracket(OpenBracket, "")
	subOut.eopt.send(p, pkt)

	for {
		pkt = subOut.ipt.receive(p)
		if pkt == nil {
			//trace(p, " Received end of stream")
			break
		}
		pkt.owner = p

		//fmt.Println(pkt.Contents)

		subOut.eopt.send(p, pkt)

	}
	pkt = p.CreateBracket(CloseBracket, "")
	subOut.eopt.send(p, pkt)
}
