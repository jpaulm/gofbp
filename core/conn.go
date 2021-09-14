package core

type Conn interface {
	receive(p *Process) *Packet
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
	GetType() string

	GetArrayItem(i int) Conn
	ArrayLength() int
}
