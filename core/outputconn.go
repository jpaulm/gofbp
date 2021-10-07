package core

type OutputConn interface {
	send(*Process, *Packet) bool

	//IsEmpty() bool
	//IsClosed() bool
	SetOptional(b bool)
	//GetType() string
	Close()
}

type OutputArrayConn interface {
	OutputConn
	GetArrayItem(i int) *OutPort
	SetArrayItem(c *OutPort, i int)
	ArrayLength() int
}
