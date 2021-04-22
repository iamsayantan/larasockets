package larasockets

// Connection interface defines the method for an individual connection to the server
type Connection interface {
	// Id returns the unique identifier for the particular
	// Connection instance
	Id() string

	// App returns the application the connection was made for
	App() *Application

	// Send will send the given data back to the client
	Send(data interface{})

	// Receive reads the data from the connection
	Receive()

	// Close closes the current connection
	Close()
}
