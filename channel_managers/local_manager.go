package channel_managers

import (
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/channels"
	"github.com/iamsayantan/larasockets/events"
	"go.uber.org/zap"
)

type localChannelManager struct {
	appManager larasockets.ApplicationManager

	logger *zap.Logger
	// channels stores all the active channels per app with channel name as the key.
	// map[appId]map[channelName]*Channel2
	channels map[string]map[string]larasockets.Channel
}

// NewLocalManager will return return a ChannelManager instance that is managed in memory.
func NewLocalManager(apps larasockets.ApplicationManager, logger *zap.Logger) larasockets.ChannelManager {
	channelManager := &localChannelManager{
		appManager: apps,
		logger:     logger,
	}

	channelManager.channels = make(map[string]map[string]larasockets.Channel, 0)
	return channelManager
}

func (cm *localChannelManager) AppManager() larasockets.ApplicationManager {
	return cm.appManager
}

func (cm *localChannelManager) FindChannel(appId, channelName string) larasockets.Channel {
	channel, ok := cm.channels[appId][channelName]
	if !ok {
		return nil
	}

	return channel
}

func (cm *localChannelManager) AllChannels(appId string) []larasockets.Channel {
	c := make([]larasockets.Channel, 0)
	existingChannels, ok := cm.channels[appId]
	if !ok {
		cm.logger.Info("no channel found for this app id", zap.String("application_id", appId))
		return c
	}

	for _, channel := range existingChannels {
		c = append(c, channel)
	}

	return c
}

func (cm *localChannelManager) ConcurrentConnectionsForApp(appId string) int {
	activeChannels := cm.AllChannels(appId)
	if len(activeChannels) == 0 {
		return 0
	}

	connectionIds := make([]string, 0)
	connectionKeys := make(map[string]bool, 0)

	for _, c := range activeChannels {
		for _, conn := range c.Connections() {
			if _, ok := connectionKeys[conn.Id()]; !ok {
				connectionIds = append(connectionIds, conn.Id())
				connectionKeys[conn.Id()] = true
			}
		}
	}

	return len(connectionIds)
}

func (cm *localChannelManager) FindOrCreateChannel(appId, channelName string) larasockets.Channel {
	existingChannel := cm.FindChannel(appId, channelName)
	if existingChannel != nil {
		return existingChannel
	}

	// determine if a channels map already exists for this app, if no channel map exists,
	// we need to make a new map for it.
	existingChannelsForApp, ok := cm.channels[appId]
	if !ok {
		existingChannelsForApp = make(map[string]larasockets.Channel)
	}

	newChannel := channels.NewChannel(channelName)
	existingChannelsForApp[channelName] = newChannel

	cm.channels[appId] = existingChannelsForApp
	return newChannel
}

func (cm *localChannelManager) SubscribeToChannel(conn larasockets.Connection, channelName string, payload interface{}) {
	channel := cm.FindOrCreateChannel(conn.App().Id(), channelName)
	channel.Subscribe(conn, payload)

	// if this is the first connection in the channel, then we can trigger a channel-occupied event.
	if currentConns := channel.Connections(); len(currentConns) == 1 {
		events.LogEvent(cm, events.Occupied, events.DashboardLogDetails{
			AppId:       conn.App().Id(),
			ChannelName: channelName,
		})
	}
}

func (cm *localChannelManager) UnsubscribeFromChannel(conn larasockets.Connection, channelName string, payload interface{}) {
	channel := cm.FindChannel(conn.App().Id(), channelName)
	if channel == nil {
		cm.logger.Error("channel not found", zap.String("application_id", conn.App().Id()), zap.String("channel_name", channelName))
		return
	}

	channel.UnSubscribe(conn)

	// if there are no  more connections is the channel, then we trigger a channel-vacated event.
	if currentConns := channel.Connections(); len(currentConns) == 0 {
		events.LogEvent(cm, events.Vacated, events.DashboardLogDetails{
			AppId:       conn.App().Id(),
			ChannelName: channelName,
		})
	}
}

// UnsubscribeFromAllChannels will unsubscribe the connection from all the channels it is subscribed to
func (cm *localChannelManager) UnsubscribeFromAllChannels(conn larasockets.Connection) {
	c := cm.AllChannels(conn.App().Id())
	for _, channel := range c {
		channel.UnSubscribe(conn)
		// if there are no  more connections is the channel, then we trigger a channel-vacated event.
		if currentConns := channel.Connections(); len(currentConns) == 0 {
			events.LogEvent(cm, events.Vacated, events.DashboardLogDetails{
				AppId:       conn.App().Id(),
				ChannelName: channel.Name(),
			})
		}
	}
}
