package events

type ConnectionEstablishedEvent struct {
	Event string                    `json:"event"`
	Data  ConnectionEstablishedData `json:"data"`
}

type ConnectionEstablishedData struct {
	SocketId        string `json:"socket_id"`
	ActivityTimeout int    `json:"activity_timeout"`
}

func NewConnectionEstablished(socketId string, activityTimeout int) *ConnectionEstablishedEvent {
	eventData := ConnectionEstablishedData{
		SocketId:        socketId,
		ActivityTimeout: activityTimeout,
	}

	event := ConnectionEstablishedEvent{
		Event: "pusher:connection_established",
		Data:  eventData,
	}

	return &event
}
