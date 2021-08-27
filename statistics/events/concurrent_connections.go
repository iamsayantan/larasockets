package events

import (
	"encoding/json"
	"fmt"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/statistics"
)

func ConcurrentConnectionChanged(cm larasockets.ChannelManager, stat statistics.Statistic) {
	channelName := fmt.Sprintf("private-app-%s-stats-concurrent-connections", stat.AppId())
	channel := cm.FindOrCreateChannel(stat.AppId(), channelName)

	data := make(map[string]int, 0)
	data["concurrent_connections"] = stat.ConcurrentConnections()

	payloadData, err := json.Marshal(data)
	if err != nil {
		return
	}

	msg := messages.PusherEventPayload{
		Event:   "update",
		Channel: channelName,
		Data:    string(payloadData),
	}

	channel.Broadcast(msg)
}
