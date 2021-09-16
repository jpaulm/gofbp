package core

//var _ Conn = (*InArrayPort)(nil)

type OutArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []*OutPort
	closed   bool
}

/*
func (c *OutArrayPort) isDrained() bool {
	for _, v := range c.array {
		if !v.isDrained() {
			return false
		}
	}
	return true
}

func (c *OutArrayPort) IsEmpty() bool {
	return false
}

func (c *OutArrayPort) receive(p *Process) *Packet {
	return nil
}

func (c *OutArrayPort) IsClosed() bool {
	return c.closed
}
*/

//func (c *OutArrayPort) resetForNextExecution() {}

func (c *OutArrayPort) GetType() string {
	return "OutArrayPort"
}

func (c *OutArrayPort) GetArrayItem(i int) *OutPort {
	if i >= len(c.array) {
		return nil
	}
	return c.array[i]
}

func (c *OutArrayPort) SetArrayItem(c2 *OutPort, i int) {
	if i >= len(c.array) {
		// add to .array to fit c2
		increaseBy := make([]*OutPort, i-len(c.array)+1)
		c.array = append(c.array, increaseBy...)

	}
	c.array[i] = c2
}

func (c *OutArrayPort) ArrayLength() int {
	return len(c.array)
}
