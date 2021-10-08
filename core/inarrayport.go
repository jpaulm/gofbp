package core

type InArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []InputConn
}

func (c *InArrayPort) isDrained() bool {
	for _, v := range c.array {
		if !v.isDrained() {
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

func (c *InArrayPort) IsClosed() bool {
	for _, v := range c.array {
		if !v.IsClosed() {
			return false
		}
	}
	return true
}

func (c *InArrayPort) resetForNextExecution() {}

func (c *InArrayPort) GetArrayItem(i int) InputConn {
	if i >= len(c.array) {
		return nil
	}
	return c.array[i]
}

func (c *InArrayPort) SetArrayItem(c2 InputConn, i int) {
	if i >= len(c.array) {
		// add to .array to fit c2
		increaseBy := make([]InputConn, i-len(c.array)+1)
		c.array = append(c.array, increaseBy...)
	}
	c.array[i] = c2
}

func (c *InArrayPort) ArrayLength() int {
	return len(c.array)
}
