package messages

// PusherSubscriptionPayload represents the payload structure for a subscription request from the client
type PusherSubscriptionPayload struct {
	// Channel is the name of the channel that is being subscribed to
	Channel string `json:"channel"`

	// Auth is an optional string. If the channel being subscribed to is a presence or private channel,
	// then the subscription needs to be authenticated. The authentication signature should be provided
	// if required. This value is generated in the application server.
	Auth string `json:"auth"`

	// ChannelData is additional optional information about the channel if the channel is a presence
	// channel. The JSON of the channel data will be generated in the application server and encoded
	// as string.
	ChannelData string `json:"channel_data"`
}

// PusherUnsubscribePayload represents the payload for channel unsubscribe request.
type PusherUnsubscribePayload struct {
	// Channel is the name of the channel being unsubscribed to
	Channel string `json:"channel"`
}

// PusherEventPayload is the payload structure for events that are received from the application servers
type PusherEventPayload struct {
	// Channel is the name of the channel where the event needs to be delivered
	Channel string `json:"channel"`

	// Event is the name of the event that is triggered
	Event string `json:"event"`

	// Data is any additional data that is sent along with the event
	Data string `json:"data"`
}
