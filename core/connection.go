package core

import (
	"fmt"
	"reflect"
	"sync"
)

// https://stackoverflow.com/questions/36857167/how-to-correctly-use-sync-cond

type Connection struct {
	network   *Network
	pktArray  []*Packet
	is, ir    int // send index and receive index
	mtx       sync.Mutex
	condNE    *sync.Cond
	condNF    *sync.Cond
	closed    bool
	UpStrmCnt int
	portName  string
	fullName  string
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.condNF.L.Lock()
	fmt.Println(p.Name, "Sending", pkt.Contents)
	for c.IsFull() { // connection is full
		c.condNF.Wait()
	}
	fmt.Println(p.Name+" Sent ", pkt.Contents)
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Broadcast()
	c.condNF.L.Unlock()
	return true
}

func (c *Connection) receive(p *Process) *Packet {
	c.condNE.L.Lock()
	fmt.Println(p.Name + " Receiving ")
	if c.IsEmpty() { // connection is empty
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

func (c *Connection) Close() {
	c.mtx.Lock()
	c.closed = true
	c.mtx.Unlock()
}

func (c *Connection) IsEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *Connection) IsClosed() bool {
	return c.closed
}

func (c *Connection) IsFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}

func (c *Connection) Lock()   { c.mtx.Lock() }
func (c *Connection) Unlock() { c.mtx.Unlock() }
