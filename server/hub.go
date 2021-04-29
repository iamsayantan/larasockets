package server

import (
	"github.com/iamsayantan/larasockets"
	"go.uber.org/zap"
)

// Hub
type Hub struct {
	logger        *zap.Logger
	channelManger larasockets.ChannelManager

	// connections holds all the active connections to our servers. Each connection is
	// assigned an unique id, so this is a map whose key is the id of the connection.
	connections map[string]larasockets.Connection

	// register channel listens for new connections. Any new connection to our server
	// would be added to the hub via this channel.
	register chan larasockets.Connection

	// unregister channel listens for connections that are being closed so that they
	// can be removed.
	unregister chan larasockets.Connection
}

// NewHub returns pointer to a new Hub instance.
func NewHub(logger *zap.Logger, cm larasockets.ChannelManager) *Hub {
	hub := &Hub{
		logger:        logger,
		channelManger: cm,
		connections:   make(map[string]larasockets.Connection, 0),
		register:      make(chan larasockets.Connection),
		unregister:    make(chan larasockets.Connection),
	}

	go hub.run()
	return hub
}

// run listens for new connections in the hub and
func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn.Id()] = conn
		case conn := <-h.unregister:
			if _, ok := h.connections[conn.Id()]; ok {
				delete(h.connections, conn.Id())
			}
		}
	}
}

// RemoveConnection will remove the connection from all the channels, and also remove it
// form the hub.
func (h *Hub) RemoveConnection(conn larasockets.Connection) {
	h.channelManger.UnsubscribeFromAllChannels(conn)
	h.unregister <- conn
}
