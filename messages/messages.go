package messages

import (
	"encoding/json"
	"github.com/iamsayantan/larasockets"
	"strings"
)

// PusherIncomingMessagePayload is the raw message that is received from the client
type PusherIncomingMessagePayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type PusherOutgoingMessagePayload struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewPusherMessage(conn larasockets.Connection, cm larasockets.ChannelManager, payload PusherIncomingMessagePayload) larasockets.PusherMessage {
	if strings.HasPrefix(payload.Event, "pusher:") {
		return newPusherProtocolMessage(conn, cm, payload)
	}

	return newPusherClientMessage(conn, cm, payload)
}

func NewPusherErrorMessage(message string, code int) *PusherOutgoingMessagePayload {
	dataPayload := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{Message: message, Code: code}

	return &PusherOutgoingMessagePayload{
		Event: "pusher:error",
		Data:  dataPayload,
	}
}
