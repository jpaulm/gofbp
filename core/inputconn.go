package core

type InputConn interface {
	receive(p *Process) *Packet
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
	//GetType() string

	GetArrayItem(i int) *Connection
	SetArrayItem(c *Connection, i int)
	ArrayLength() int
}
