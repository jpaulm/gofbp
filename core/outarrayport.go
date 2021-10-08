package core

type OutArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []OutputConn
	closed   bool
}

func (o *OutArrayPort) GetArrayItem(i int) OutputConn {
	if i >= len(o.array) {
		return nil
	}
	return o.array[i]
}

func (o *OutArrayPort) SetArrayItem(o2 OutputConn, i int) {
	if i >= len(o.array) {
		// add to .array to fit c2
		increaseBy := make([]OutputConn, i-len(o.array)+1)
		o.array = append(o.array, increaseBy...)
	}
	o.array[i] = o2
}

func (o *OutArrayPort) ArrayLength() int {
	return len(o.array)
}

func (o *OutArrayPort) Close() {
	for _, v := range o.array {
		v.Close()
	}
}
