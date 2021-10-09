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
	outputCommon
	GetArrayItem(i int) OutputConn
	SetArrayItem(c OutputConn, i int)
	ArrayLength() int
}
