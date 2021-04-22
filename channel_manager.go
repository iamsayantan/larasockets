package larasockets

// ChannelManager interface defines methods to manage channels in the server.
type ChannelManager interface {
	// AppManager returns the AppManager instance, AppManager manages all the applications in the system
	AppManager() ApplicationManager

	// FindChannel finds a channel. Each channel should be stored per app basis. So it would be possible
	// to multiple app have channels with same name
	FindChannel(appId, channelName string) Channel

	// AllChannels will return all the channels currently in the manager.
	AllChannels(appId string) []Channel

	// FindOrCreateChannel creates a new Channel for the application if it doesn't find it.
	FindOrCreateChannel(appId, channelName string) Channel

	// SubscribeToChannel subscribes a connection to a channel.
	SubscribeToChannel(conn Connection, channelName string, payload interface{})

	// UnsubscribeFromChannel removes a connection from the channel
	UnsubscribeFromChannel(conn Connection, channelName string, payload interface{})

	// UnsubscribeFromAllChannels will unsubscribe the connection across all the channels it is subscribed to
	UnsubscribeFromAllChannels(conn Connection)
}
