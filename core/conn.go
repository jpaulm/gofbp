package core

type Conn interface {
	receive(p *Process) *Packet

	IsEmpty() bool
	IsClosed() bool
	ResetClosed()
	GetType() string
}
