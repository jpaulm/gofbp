package core

type InputConn interface {
	receive(p *Process) *Packet
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
	//GetType() string
}

type InputArrayConn interface {
	InputConn
	GetArrayItem(i int) *Connection
	SetArrayItem(c *Connection, i int)
	ArrayLength() int
}
