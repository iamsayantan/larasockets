package events

import (
	"encoding/json"
	"fmt"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"time"
)

type EventType string

const (
	Connected        EventType = "connected"
	Disconnected               = "disconnected"
	Occupied                   = "occupied"
	Vacated                    = "vacated"
	Subscribed                 = "subscribed"
	WebsocketMessage           = "websocket_message"
	ApiMessage                 = "api_message"
)

type DashboardLogDetails struct {
	AppId        string
	ChannelName  string
	EventName    string
	ConnectionId string
	EventPayload string
}

func LogEvent(cm larasockets.ChannelManager, eventType EventType, details DashboardLogDetails) {
	if details.AppId == "" {
		return
	}

	dashboardLogChannelName := fmt.Sprintf("private-websockets-dashboard-%s", details.AppId)
	logChannel := cm.FindOrCreateChannel(details.AppId, dashboardLogChannelName)

	logEventPayload := struct {
		Type         EventType `json:"type"`
		Time         int64     `json:"time"`
		EventName    string    `json:"event_name"`
		ChannelName  string    `json:"channel_name"`
		ConnectionId string    `json:"connection_id"`
		Payload      string    `json:"payload"`
	}{
		Type:         eventType,
		Time:         time.Now().Unix(),
		EventName:    details.EventName,
		ChannelName:  details.ChannelName,
		ConnectionId: details.ConnectionId,
		Payload:      details.EventPayload,
	}

	payloadData, err := json.Marshal(logEventPayload)
	if err != nil {
		return
	}

	msg := messages.PusherEventPayload{
		Event:   "log",
		Channel: dashboardLogChannelName,
		Data:    string(payloadData),
	}

	logChannel.Broadcast(msg)
}
