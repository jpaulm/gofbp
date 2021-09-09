package core

type Conn interface {
	Receive() *Packet
}
