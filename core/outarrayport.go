package core

//var _ Conn = (*InArrayPort)(nil)

type OutArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []*OutPort
	closed   bool
}

func (c *OutArrayPort) SetOptional(b bool) {}

func (c *OutArrayPort) GetType() string {
	return "OutArrayPort"
}

func (c *OutArrayPort) GetArrayItem(i int) *OutPort {
	if i >= len(c.array) {
		return nil
	}
	return c.array[i]
}

func (o *OutArrayPort) SetArrayItem(o2 *OutPort, i int) {
	if i >= len(o.array) {
		// add to .array to fit c2
		increaseBy := make([]*OutPort, i-len(o.array)+1)
		o.array = append(o.array, increaseBy...)
	}
	o.array[i] = o2
}

func (c *OutArrayPort) ArrayLength() int {
	return len(c.array)
}
