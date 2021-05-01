package dto

type DailyStatSnapshot struct {
	PeakConnections   int `json:"peak_connections"`
	ApiMessages       int `json:"api_messages"`
	WebsocketMessages int `json:"websocket_messages"`
}
