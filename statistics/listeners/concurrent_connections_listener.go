package listeners

import (
	"encoding/json"
	"fmt"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/statistics"
)

func NewConcurrentConnectionListener(cm larasockets.ChannelManager) statistics.StatsCollectionListener {
	return &concurrentConnectionListener{channelManager: cm}
}

type concurrentConnectionListener struct {
	channelManager larasockets.ChannelManager
}

func (c *concurrentConnectionListener) ListenStatChanged(stat statistics.Statistic) {
	channelName := fmt.Sprintf("private-app-%s-stats-concurrent-connections", stat.AppId())
	channel := c.channelManager.FindOrCreateChannel(stat.AppId(), channelName)

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
