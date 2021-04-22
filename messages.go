package larasockets

// PusherMessage interface defines a single method Respond which must be implemented
// by all the messages coming from client.
type PusherMessage interface {
	Respond()
}
