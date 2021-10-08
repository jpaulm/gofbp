package core

type NullOutPort struct {
	name string
}

func (c *NullOutPort) SetOptional(b bool) {}

func (c *NullOutPort) send(*Process, *Packet) bool { panic("send on null port") }

func (c *NullOutPort) Close() {}
