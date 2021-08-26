package stores

import (
	"github.com/iamsayantan/larasockets/statistics"
	"gorm.io/gorm"
	"time"
)

type LarasocketsStatistic struct {
	ID                uint   `json:"id" gorm:"primarykey"`
	AppId             string `gorm:"index"`
	PeakConnections   int
	WebsocketMessages int
	ApiMessages       int
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

func NewDatabaseStorage(db *gorm.DB) statistics.StatsStorage {
	return &dbStore{db: db}
}

type dbStore struct {
	db *gorm.DB
}

func (m *dbStore) Store(statistic statistics.Statistic) {
	statToStore := LarasocketsStatistic{
		AppId:             statistic.AppId(),
		PeakConnections:   statistic.PeakConnections(),
		WebsocketMessages: statistic.WebsocketMessages(),
		ApiMessages:       statistic.ApiMessages(),
	}

	m.db.Create(&statToStore)
}

func (m *dbStore) DailyStatForApp(appId string) *statistics.Statistic {
	var stats LarasocketsStatistic
	err := m.db.Select("MAX(peak_connections) AS peak_connections, SUM(websocket_messages) AS websocket_messages, SUM(api_messages) AS api_messages").
		Where("app_id = ?", appId).
		Where("created_at >= CURDATE()").
		Where("created_at < CURDATE() + INTERVAL 1 DAY").
		First(&stats).
		Error

	if err != nil {
		return nil
	}

	return statistics.NewStatisticWithData(appId, 0, stats.PeakConnections, stats.WebsocketMessages, stats.ApiMessages)
}

func (m *dbStore) StatsByTimeRange(appId string, startTime time.Time, endTime time.Time) *statistics.StatisticByTime {
	var stats []LarasocketsStatistic
	statsResponse := statistics.NewStatisticByTime()

	m.db.Where("app_id = ?", appId).Where("created_at >= ?", startTime).Where("created_at <= ?", endTime).Order("created_at DESC").Find(&stats)
	for _, stat := range stats {
		statsResponse.Set(stat.CreatedAt, statistics.NewStatisticWithData(appId, 0, stat.PeakConnections, stat.WebsocketMessages, stat.ApiMessages))
	}

	return statsResponse
}
