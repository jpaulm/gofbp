package core

// https://stackoverflow.com/questions/36857167/how-to-correctly-use-sync-cond

type InitializationConnection struct {
	network  *Network
	portName string
	fullName string
}

func (c *InitializationConnection) Lock()          { /* TODO */ }
func (c *InitializationConnection) Unlock()        { /* TODO */ }
func (c *InitializationConnection) IsEmpty() bool  { /* TODO */ return true }
func (c *InitializationConnection) IsClosed() bool { /* TODO */ return true }

func (c *InitializationConnection) receive(p *Process) *Packet {
	/* TODO */

	/*
		c.condNE.L.Lock()
		fmt.Println(p.Name + " Receiving ")
		if c.IsEmpty() { // connection is empty
			if c.closed {
				c.condNF.Broadcast()
				c.condNE.L.Unlock()
				return nil
			}
			c.condNE.Wait()
		}
		pkt := c.pktArray[c.ir]
		c.pktArray[c.ir] = nil
		fmt.Println(p.Name, "Received", pkt.Contents)
		c.ir = (c.ir + 1) % len(c.pktArray)
		pkt.owner = p
		p.ownedPkts++
		c.condNF.Broadcast()
		c.condNE.L.Unlock()
	*/
	return nil
}

/*
func (c *Connection) Close() {
	c.mtx.Lock()
	c.closed = true
	c.mtx.Unlock()
}

func (c *Connection) IsEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *Connection) IsFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}
*/
