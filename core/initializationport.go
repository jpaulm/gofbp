package core

import (
	"fmt"
	"sync"
)

type InitializationPort struct {
	network  *Network
	portName string
	fullName string
	closed   bool
	value    string
	mtx      sync.Mutex
}

func (c *InitializationPort) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationPort) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return !c.closed
}

func (c *InitializationPort) receive(p *Process) *Packet {

	if c.closed {
		return nil
	}
	fmt.Println(p.name, "Receiving IIP")
	var pkt *Packet = new(Packet)
	pkt.Contents = c.value
	pkt.owner = p
	p.ownedPkts++
	c.Close()
	fmt.Println(p.name, "Received IIP: ", pkt.Contents)
	return pkt
}

func (c *InitializationPort) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = true

}

func (c *InitializationPort) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationPort) resetForNextExecution() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = false
}

//func (c *InitializationPort) GetType() string {
//	return "InitializationPort"
//}

//func (c *InitializationPort) GetArrayItem(i int) *InPort {
//	return nil
//}

//func (c *InitializationPort) SetArrayItem(c2 *InPort, i int) {}

//func (c *InitializationPort) ArrayLength() int {
//	return 0
//}
