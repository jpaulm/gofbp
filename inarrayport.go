package gofbp

import (
	//"fmt"
	"sync"
)

type InArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []*InPort
	mtx      sync.Mutex
	closed   bool
}

func (c *InArrayPort) IsDrained() bool {
	for _, v := range c.array {
		if !v.IsDrained() {
			return false
		}
	}
	return true
}

func (c *InArrayPort) IsEmpty() bool {
	for _, v := range c.array {
		if !v.IsEmpty() {
			return false
		}
	}
	return true
}

func (c *InArrayPort) receive(p *Process) *Packet {
	panic("receive from an array port")
}

func (c *InArrayPort) IsClosed() bool {
	for _, v := range c.array {
		if !v.IsClosed() {
			return false
		}
	}
	return true
}

func (c *InArrayPort) resetForNextExecution() {}

func (c *InArrayPort) GetArrayItem(i int) *InPort {
	if i >= len(c.array) {
		return nil
	}
	return c.array[i]
}

func (c *InArrayPort) SetArrayItem(c2 *InPort, i int) {
	if i >= len(c.array) {
		// add to .array to fit c2
		increaseBy := make([]*InPort, i-len(c.array)+1)
		c.array = append(c.array, increaseBy...)
	}
	c.array[i] = c2
}

func (c *InArrayPort) ArrayLength() int {
	return len(c.array)
}

func (c *InArrayPort) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.closed = true
	for _, w := range c.array {
		w.closed = true
		BdcastTr(w.condNE, "bdcast in array NE", w.downStrProc)
	}
}
