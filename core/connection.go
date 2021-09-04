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
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.condNF.L.Lock()
	v := reflect.ValueOf(pkt.Contents) // display contents - assume string
	s := v.String()
	fmt.Println(p.Name + " Sending " + s)
	for (c.ir == c.is) && (c.pktArray[c.is] != nil) { // connection is full
		c.condNF.Wait()
	}
	fmt.Println(p.Name + " Sent " + s)
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Broadcast()
	c.condNF.L.Unlock()
	return true
}

func (p *Process) Receive(c *Connection) *Packet {

	c.condNE.L.Lock()
	fmt.Println(p.Name + " Receiving ")
	if (c.ir == c.is) && (c.pktArray[c.is] == nil) { // connection is empty
		if c.closed {
			c.condNF.Broadcast()
			c.condNE.L.Unlock()
			return nil
		}
		c.condNE.Wait()
	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	v := reflect.ValueOf(pkt.Contents) // display contents - assume string
	s := v.String()
	fmt.Println(p.Name + " Received " + s)
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Broadcast()
	c.condNE.L.Unlock()
	return pkt
}
