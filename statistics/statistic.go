package statistics

type Statistic struct {
	appId                  string
	concurrentConnections  int
	peakConnections        int
	websocketMessagesCount int
	apiMessagesCount       int
}

func NewStatistic(appId string) *Statistic {
	return &Statistic{appId: appId}
}

/*
func (s *Statistic) SetConcurrentConnection(connectionCount int) {
	s.concurrentConnections = connectionCount
}

func (s *Statistic) SetPeakConnections(peakConnectionCount int) {
	s.peakConnections = peakConnectionCount
}

func (s *Statistic) SetWebsocketMessageCount(websocketMessageCount int) {
	s.websocketMessagesCount = websocketMessageCount
}

func (s *Statistic) SetApiMessageCount(apiMessageCount int) {
	s.apiMessagesCount = apiMessageCount
}
*/
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

func (s *Statistic) Reset() {
	s2 := NewStatistic(s.appId)
	*s = *s2
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
