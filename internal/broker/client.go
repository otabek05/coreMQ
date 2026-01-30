package broker

import "rawMQ/common"

func (b *Broker) addClient(c *common.Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[c.ID] = c
}


func (b *Broker) getClient(id string) (*common.Client, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	c, ok := b.clients[id]
	return c, ok
}


func (b *Broker) removeClient(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, id)
}
