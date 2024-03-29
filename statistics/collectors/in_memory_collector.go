package collectors

import (
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/statistics"
	"github.com/iamsayantan/larasockets/statistics/events"
	"time"
)

// NewMemoryCollector returns a new stats collector that stores the data in memory.
func NewMemoryCollector(cm larasockets.ChannelManager, store statistics.StatsStorage) statistics.StatsCollector {
	collector := &memoryCollector{
		stats:     make(map[string]*statistics.Statistic),
		listeners: make([]statistics.StatsCollectionListener, 0),
		store:     store,
		cm:        cm,
	}

	go collector.periodicDumpToStorage()

	return collector
}

type memoryCollector struct {
	stats     map[string]*statistics.Statistic
	listeners []statistics.StatsCollectionListener
	store     statistics.StatsStorage
	cm        larasockets.ChannelManager
}

func (c *memoryCollector) DumpToStorage(store statistics.StatsStorage) {
	for _, stat := range c.stats {
		if stat.CanRemoveStatistics() {
			delete(c.stats, stat.AppId())
			return
		}

		store.Store(*stat)
		events.StatisticsUpdated(c.cm, *stat)

		concurrentConnections := c.cm.ConcurrentConnectionsForApp(stat.AppId())
		stat.Reset(concurrentConnections)
	}
}

func (c *memoryCollector) HandleWebsocketMessage(appId string) {
	c.findOrMake(appId).HandleNewWebsocketMessage()
}

func (c *memoryCollector) HandleApiMessage(appId string) {
	c.findOrMake(appId).HandleNewApiMessage()
}

func (c *memoryCollector) HandleConnection(appId string) {
	c.findOrMake(appId).HandleNewConnection()
	events.ConcurrentConnectionChanged(c.cm, c.GetAppStatistics(appId))
}

func (c *memoryCollector) HandleDisconnection(appId string) {
	c.findOrMake(appId).HandleDisconnection()
	events.ConcurrentConnectionChanged(c.cm, c.GetAppStatistics(appId))
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

func (c *memoryCollector) periodicDumpToStorage() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			c.DumpToStorage(c.store)
		}
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
