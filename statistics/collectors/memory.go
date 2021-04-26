package collectors

import (
	"github.com/iamsayantan/larasockets/statistics"
)

// NewMemoryCollector returns a new stats collector that stores the data in memory.
func NewMemoryCollector() statistics.StatsCollector {
	return &memoryCollector{
		stats:     make(map[string]*statistics.Statistic),
		listeners: make([]statistics.StatsCollectionListener, 0),
	}
}

type memoryCollector struct {
	stats     map[string]*statistics.Statistic
	listeners []statistics.StatsCollectionListener
}

func (c *memoryCollector) HandleWebsocketMessage(appId string) {
	c.findOrMake(appId).HandleNewWebsocketMessage()
	c.sendUpdatedStatToListeners(appId)
}

func (c *memoryCollector) HandleApiMessage(appId string) {
	c.findOrMake(appId).HandleNewApiMessage()
	c.sendUpdatedStatToListeners(appId)
}

func (c *memoryCollector) HandleConnection(appId string) {
	c.findOrMake(appId).HandleNewConnection()
	c.sendUpdatedStatToListeners(appId)
}

func (c *memoryCollector) HandleDisconnection(appId string) {
	c.findOrMake(appId).HandleDisconnection()
	c.sendUpdatedStatToListeners(appId)
}

func (c *memoryCollector) Flush() {
	c.stats = make(map[string]*statistics.Statistic, 0)
}

func (c *memoryCollector) GetAllStatistics() []statistics.Statistic {
	stats := make([]statistics.Statistic, 0)
	for _, stat := range c.stats {
		stats = append(stats, *stat)
	}

	return stats
}

func (c *memoryCollector) GetAppStatistics(appId string) statistics.Statistic {
	stat := c.findOrMake(appId)
	return *stat
}

func (c *memoryCollector) RegisterStatsListener(listener statistics.StatsCollectionListener) {
	c.listeners = append(c.listeners, listener)
}

func (c *memoryCollector) sendUpdatedStatToListeners(appId string) {
	stats := c.findOrMake(appId)
	for _, listener := range c.listeners {
		listener.ListenStatChanged(*stats)
	}
}

func (c *memoryCollector) findOrMake(appId string) *statistics.Statistic {
	if stat, ok := c.stats[appId]; ok {
		return stat
	}

	stat := statistics.NewStatistic(appId)
	c.stats[appId] = stat

	return stat
}
