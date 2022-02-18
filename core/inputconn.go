package core

type inputCommon interface {
	isDrained() bool
	resetForNextExecution()
	isEmpty() bool
	isClosed() bool
	Close()
	receive(*Process) *Packet
	//Name() string
	//SetName(string)
}

type InputConn interface {
	inputCommon
	//receive(p *Process) *Packet
}

type InputArrayConn interface {
	inputCommon

	GetArrayItem(i int) *InPort
	setArrayItem(c *InPort, i int)
	ArrayLength() int
	GetArray() []*InPort
}
