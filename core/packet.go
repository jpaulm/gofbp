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
	chains   map[string]*ChainHdr
	next     *Packet
}

type ChainHdr struct {
	name  string
	first *Packet
	last  *Packet
}
