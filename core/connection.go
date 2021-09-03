package core

import (
	"fmt"
	"reflect"
	"sync"
)

type Connection struct {
	//procs   map[string]Process
	network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	mtx      sync.Mutex
	pktArray []Packet
	is, ir   int
}

func (p *Process) Send(c *Connection, pkt *Packet) bool {
	c.mtx.Lock()
	//s, ok := pkt.contents(.string)
	v := reflect.ValueOf(pkt.contents)
	s := v.String()
	fmt.Println("Sending " + s)
	c.pktArray[c.is] = *pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	c.mtx.Unlock()
	return true
}
