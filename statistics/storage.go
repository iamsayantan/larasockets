package statistics

// StatsStorage interface defines methods for storing and fetching
// the statistics data.
type StatsStorage interface {
	Store(statistic Statistic)
}
