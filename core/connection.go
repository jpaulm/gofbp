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
	mtx    sync.Mutex
	slice  []Packet
	is, ir int
}

func (p *Process) Send(c *Connection, pkt *Packet) bool {
	c.mtx.Lock()
	//s, ok := pkt.contents(.string)
	v := reflect.ValueOf(pkt.contents)
	s := v.String()
	fmt.Println("Sending " + s)
	c.slice[c.is] = *pkt
	c.is = (c.is + 1) % len(c.slice)
	pkt.owner = nil
	c.mtx.Unlock()
	return true
}
