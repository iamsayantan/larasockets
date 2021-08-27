package dto

type DailyStatSnapshot struct {
	ConcurrentConnection int `json:"concurrent_connection"`
	PeakConnections      int `json:"peak_connections"`
	ApiMessages          int `json:"api_messages"`
	WebsocketMessages    int `json:"websocket_messages"`
}

type StatisticsPlot struct {
	X []int64 `json:"x"` // array of unix timestamps
	Y []int   `json:"y"` // value for the corresponding timestamp
}

type StatisticsGraph struct {
	ApiStats              StatisticsPlot `json:"api_stats"`
	PeakConnectionStats   StatisticsPlot `json:"peak_connection_stats"`
	WebsocketMessageStats StatisticsPlot `json:"websocket_message_stats"`
}
