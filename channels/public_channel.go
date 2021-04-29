package channels

import (
	"github.com/iamsayantan/larasockets"
)

type publicChannel struct {
	name        string
	connections map[string]larasockets.Connection
}

func (c *publicChannel) Name() string {
	return c.name
}

func (c *publicChannel) Connections() []larasockets.Connection {
	connections := make([]larasockets.Connection, 0)
	for _, conn := range c.connections {
		connections = append(connections, conn)
	}

	return connections
}

func (c *publicChannel) Subscribe(conn larasockets.Connection, payload interface{}) {
	if c.IsSubscribed(conn) {
		return
	}
	c.connections[conn.Id()] = conn

	resp := struct {
		Event   string      `json:"event"`
		Channel string      `json:"channel"`
		Data    interface{} `json:"data"`
	}{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: c.Name(),
		Data:    "{}",
	}
	conn.Send(resp)
}

func (c *publicChannel) UnSubscribe(conn larasockets.Connection) {
	if !c.IsSubscribed(conn) {
		return
	}

	delete(c.connections, conn.Id())
}

func (c *publicChannel) IsSubscribed(conn larasockets.Connection) bool {
	_, ok := c.connections[conn.Id()]
	return ok
}

func (c *publicChannel) Broadcast(data interface{}) {
	for _, conn := range c.connections {
		conn.Send(data)
	}
}

func (c *publicChannel) BroadcastExcept(data interface{}, excludedConnectionId string) {
	for _, conn := range c.connections {
		if conn.Id() == excludedConnectionId {
			continue
		}

		conn.Send(data)
	}
}

func newPublicChannel(name string) larasockets.Channel {
	return &publicChannel{
		name:        name,
		connections: make(map[string]larasockets.Connection, 0),
	}
}
