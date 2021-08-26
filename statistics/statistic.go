package statistics

import "time"

type Statistic struct {
	appId                  string
	concurrentConnections  int
	peakConnections        int
	websocketMessagesCount int
	apiMessagesCount       int
}

type StatisticByTime struct {
	timestamps []int64
	statistics map[int64]*Statistic
}

func NewStatistic(appId string) *Statistic {
	return &Statistic{appId: appId}
}

func NewStatisticWithData(appId string, concurrentConnections, peakConnections, websocketMessages, apiMessages int) *Statistic {
	return &Statistic{
		appId:                  appId,
		concurrentConnections:  concurrentConnections,
		peakConnections:        peakConnections,
		websocketMessagesCount: websocketMessages,
		apiMessagesCount:       apiMessages,
	}
}

func NewStatisticByTime() *StatisticByTime {
	return &StatisticByTime{
		timestamps: make([]int64, 0),
		statistics: make(map[int64]*Statistic, 0),
	}
}

func (s *Statistic) AppId() string {
	return s.appId
}

func (s *Statistic) ConcurrentConnections() int {
	return s.concurrentConnections
}

func (s *Statistic) PeakConnections() int {
	return s.peakConnections
}

func (s *Statistic) WebsocketMessages() int {
	return s.websocketMessagesCount
}

func (s *Statistic) ApiMessages() int {
	return s.apiMessagesCount
}

func (s *Statistic) HandleNewConnection() {
	s.concurrentConnections++

	if s.concurrentConnections > s.peakConnections {
		s.peakConnections = s.concurrentConnections
	}
}

func (s *Statistic) HandleDisconnection() {
	s.concurrentConnections--
}

func (s *Statistic) HandleNewWebsocketMessage() {
	s.websocketMessagesCount++
}

func (s *Statistic) HandleNewApiMessage() {
	s.apiMessagesCount++
}

func (s *Statistic) Reset(concurrentConnections int) {
	s.concurrentConnections = concurrentConnections
	s.apiMessagesCount = 0
	s.websocketMessagesCount = 0

	if concurrentConnections > 0 {
		s.peakConnections = concurrentConnections
	} else {
		s.peakConnections = 0
	}
}

// CanRemoveStatistics checks if there is any current activity in for the app,
// if there is no current activity then, this statistics can be safely removed.
func (s *Statistic) CanRemoveStatistics() bool {
	return s.peakConnections == 0 && s.concurrentConnections == 0
}

func (s *Statistic) GetCurrentSnapshot() map[string]interface{} {
	stats := make(map[string]interface{}, 0)
	stats["app_id"] = s.appId
	stats["concurrent_connections"] = s.concurrentConnections
	stats["peak_connections"] = s.peakConnections
	stats["websocket_messages"] = s.websocketMessagesCount
	stats["api_messages"] = s.apiMessagesCount

	return stats
}

func (st *StatisticByTime) Set(t time.Time, statistic *Statistic) {
	st.timestamps = append(st.timestamps, t.Unix())
	st.statistics[t.Unix()] = statistic
}

func (st *StatisticByTime) Timestamps() []int64 {
	return st.timestamps
}

func (st *StatisticByTime) Get(timestamp int64) *Statistic {
	if stat, ok := st.statistics[timestamp]; ok {
		return stat
	}

	return nil
}
