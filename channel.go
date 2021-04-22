package larasockets

// Channel interface defines the methods required for a channel to implement.
type Channel interface {
	// Name returns the name of the channel
	Name() string

	// Connections returns all the concurrent connections to this channel
	Connections() []Connection

	// Subscribe subscribes a new connection to the channel
	Subscribe(conn Connection, payload interface{})

	// Unsubscribe removes the connection from the current subscribed connections
	UnSubscribe(conn Connection)

	// IsSubscribed returns if the given connection is already subscribed to this channel
	IsSubscribed(conn Connection) bool

	// Broadcast sends the data to all the connected connections
	Broadcast(data interface{})

	// BroadcastExcept sends the data to all the connected connections except the given
	// connection. Usually its the connection who initiated the broadcast.
	BroadcastExcept(data interface{}, excludedConnectionId string)
}
