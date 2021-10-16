package core

type inputCommon interface {
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
	//GetType() string
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
