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
