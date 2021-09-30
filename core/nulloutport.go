package core

type NullOutPort struct {
	name     string
	Conn     *Connection
	optional bool
}

func (c *NullOutPort) SetOptional(b bool) {}

func (c *NullOutPort) send(*Process, *Packet) bool { panic("send on null port") }

func (c *NullOutPort) GetType() string {
	return "NullOutPort"
}

func (c *NullOutPort) GetArrayItem(i int) *OutPort {
	return nil
}

func (c *NullOutPort) SetArrayItem(o *OutPort, i int) {}

func (c *NullOutPort) ArrayLength() int {
	return 0
}
