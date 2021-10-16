package core

type outputCommon interface {
	Close()
}

type OutputConn interface {
	outputCommon
	send(*Process, *Packet) bool

	IsConnected() bool
}

type OutputArrayConn interface {
	//OutputConn
	outputCommon
	GetArrayItem(i int) *OutPort
	SetArrayItem(c *OutPort, i int)
	ArrayLength() int
}
