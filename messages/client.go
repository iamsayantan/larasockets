package messages

import "github.com/iamsayantan/larasockets"

type pusherChannelClientMessage struct {
	connection     larasockets.Connection
	channelManager larasockets.ChannelManager
	payload        PusherIncomingMessagePayload
}

func (p *pusherChannelClientMessage) Respond() {
	//panic("implement me")
}

func newPusherClientMessage(conn larasockets.Connection, cm larasockets.ChannelManager, payload PusherIncomingMessagePayload) larasockets.PusherMessage {
	pm := &pusherChannelClientMessage{
		connection:     conn,
		channelManager: cm,
		payload:        payload,
	}

	return pm
}
