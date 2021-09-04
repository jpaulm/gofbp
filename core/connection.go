package core

import (
	"fmt"
	"reflect"
	"sync"
)

// https://stackoverflow.com/questions/36857167/how-to-correctly-use-sync-cond

type Connection struct {
	network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	pktArray []*Packet
	is, ir   int // send index and receive index
	mtx      sync.Mutex
	condNE   *sync.Cond
	condNF   *sync.Cond
	closed   bool
	//l        sync.Locker
}

func (p *Process) Send(c *Connection, pkt *Packet) bool {
	c.condNF.L.Lock()
	v := reflect.ValueOf(pkt.Contents) // display contents - assume string
	s := v.String()
	fmt.Println("Sending " + s)
	for (c.ir == c.is) && (c.pktArray[c.is] != nil) { // connection is full
		c.condNF.Wait()
	}
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Broadcast()
	c.condNF.L.Unlock()
	return true
}

func (p *Process) Receive(c *Connection) *Packet {

	if (c.ir == c.is) && (c.pktArray[c.is] == nil) && c.closed {
		return nil
	}

	c.condNE.L.Lock()
	//v := reflect.ValueOf(pkt.contents)  // display contents - assume string
	//s := v.String()
	//fmt.Println("Sending " + s)
	if (c.ir == c.is) && (c.pktArray[c.is] == nil) { // connection is empty
		c.condNE.Wait()
	}
	pkt := c.pktArray[c.ir]
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Broadcast()
	c.condNE.L.Unlock()
	return pkt
}
