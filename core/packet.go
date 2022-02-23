package core

const (
	NormalPacket int32 = iota
	OpenBracket
	CloseBracket
	Signal
)

type Packet struct {
	Contents interface{}
	PktType  int32
	owner    interface{} // must be *Process or *Packet
	chains   map[string]*Chain
	next     *Packet
}

type Chain struct {
	name  string
	owner *Packet
	first *Packet
	last  *Packet
}
