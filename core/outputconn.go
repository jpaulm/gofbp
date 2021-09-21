package core

type OutputConn interface {
	//send(p *Process) *Packet

	//IsEmpty() bool
	//IsClosed() bool
	SetOptional(b bool)
	GetType() string

	GetArrayItem(i int) *OutPort
	SetArrayItem(c *OutPort, i int)
	ArrayLength() int
}
