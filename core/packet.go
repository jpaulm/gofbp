package core

const (
	Normal int32 = iota
	OpenBracket
	CloseBracket
)

type Packet struct {
	Contents interface{}
	PktType  int32
	owner    *Process
}
