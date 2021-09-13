package core

type Conn interface {
	receive(p *Process) *Packet
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
	GetType() string

	GetArray() []*Connection
	SetArray(*Connection, int)
}
