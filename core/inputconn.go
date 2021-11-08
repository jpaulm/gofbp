package core

type inputCommon interface {
	IsDrained() bool
	resetForNextExecution()
	IsEmpty() bool
	IsClosed() bool
	Close()
}

type InputConn interface {
	inputCommon
	receive(p *Process) *Packet
}

type InputArrayConn interface {
	inputCommon

	GetArrayItem(i int) *InPort
	SetArrayItem(c *InPort, i int)
	ArrayLength() int
}
