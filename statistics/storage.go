package statistics

import "time"

// StatsStorage interface defines methods for storing and fetching
// the statistics data.
type StatsStorage interface {
	Store(statistic Statistic)
	DailyStatForApp(appId string) *Statistic
	StatsByTimeRange(appId string, startTime time.Time, endTime time.Time) *StatisticByTime
}
