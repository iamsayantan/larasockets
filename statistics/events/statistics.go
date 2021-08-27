package events

import (
	"encoding/json"
	"fmt"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/statistics"
	"time"
)

func StatisticsUpdated(cm larasockets.ChannelManager, stat statistics.Statistic) {
	channelName := fmt.Sprintf("private-app-%s-current-stats", stat.AppId())
	channel := cm.FindOrCreateChannel(stat.AppId(), channelName)

	data := stat.GetCurrentSnapshot()
	data["timestamp"] = time.Now().Unix()

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
