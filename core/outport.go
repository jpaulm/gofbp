package core

type OutPort struct {
	name     string
	Conn     *Connection
	optional bool
}

func (c *OutPort) SetOptional(b bool) {
	c.optional = b
}

func (c *OutPort) GetType() string {
	return "OutPort"
}

func (c *OutPort) GetArrayItem(i int) *OutPort {
	return nil
}

func (c *OutPort) SetArrayItem(o *OutPort, i int) {}

func (c *OutPort) ArrayLength() int {
	return 0
}
