package core

type inputCommon interface {
	isDrained() bool
	resetForNextExecution()

	IsEmpty() bool
	IsClosed() bool
}

type InputConn interface {
	inputCommon
	receive(p *Process) *Packet
}

type InputArrayConn interface {
	inputCommon

	GetArrayItem(i int) InputConn
	SetArrayItem(c InputConn, i int)
	ArrayLength() int
}
