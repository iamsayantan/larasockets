package messages

import (
	"encoding/json"
	"github.com/iamsayantan/larasockets"
	"strings"
)

type pusherChannelProtocolMessage struct {
	pusherChannelClientMessage
}

func (p *pusherChannelProtocolMessage) Respond() {
	payload := p.payload

	events := strings.SplitAfter(payload.Event, ":")
	eventName := strings.Join(events[1:], "")

	switch eventName {
	case "subscribe":
		p.handleSubscription()
	case "unsubscribe":
		p.handleUnSubscribe()
	case "ping":
		p.handlePing()
	}
}

func (p *pusherChannelProtocolMessage) handleSubscription() {
	var payload PusherSubscriptionPayload
	err := json.Unmarshal(p.payload.Data, &payload)
	if err != nil {
		return
	}

	p.channelManager.SubscribeToChannel(p.connection, payload.Channel, payload)
}

func (p *pusherChannelProtocolMessage) handleUnSubscribe() {
	var payload PusherUnsubscribePayload
	if err := json.Unmarshal(p.payload.Data, &payload); err != nil {
		return
	}

	p.channelManager.UnsubscribeFromChannel(p.connection, payload.Channel, payload)
}

func (p *pusherChannelProtocolMessage) handlePing() {
	resp := struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}{
		Event: "pusher:pong",
	}

	p.connection.Send(resp)
}

func newPusherProtocolMessage(conn larasockets.Connection, cm larasockets.ChannelManager, payload PusherIncomingMessagePayload) larasockets.PusherMessage {
	return &pusherChannelProtocolMessage{pusherChannelClientMessage{
		connection:     conn,
		channelManager: cm,
		payload:        payload,
	}}
}
