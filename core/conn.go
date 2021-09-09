package core

type Conn interface {
	receive(p *Process) *Packet

	Lock()
	Unlock()

	IsEmpty() bool
	IsClosed() bool
}
