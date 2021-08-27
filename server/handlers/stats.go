package handlers

import (
	"github.com/go-chi/chi"
	"github.com/iamsayantan/larasockets/server/handlers/dto"
	"github.com/iamsayantan/larasockets/server/rendering"
	"github.com/iamsayantan/larasockets/statistics"
	"net/http"
	"time"
)

func NewStatsHandler(store statistics.StatsStorage, collector statistics.StatsCollector) *StatsHandler {
	return &StatsHandler{statsStore: store, statsCollector: collector}
}

type StatsHandler struct {
	statsStore     statistics.StatsStorage
	statsCollector statistics.StatsCollector
}

func (h *StatsHandler) GetStatForToday(w http.ResponseWriter, r *http.Request) {
	stat := h.statsStore.DailyStatForApp(chi.URLParam(r, "appId"))
	currentStat := h.statsCollector.GetAppStatistics(chi.URLParam(r, "appId"))
	resp := dto.DailyStatSnapshot{
		ConcurrentConnection: currentStat.ConcurrentConnections(),
		PeakConnections:      stat.PeakConnections(),
		ApiMessages:          stat.ApiMessages(),
		WebsocketMessages:    stat.WebsocketMessages(),
	}

	rendering.RenderSuccessWithData(w, "success", http.StatusOK, resp)
}

func (h *StatsHandler) GetStatsForGraph(w http.ResponseWriter, r *http.Request) {
	endTime := time.Now()
	startTime := endTime.Add(-time.Minute * 30)

	stats := h.statsStore.StatsByTimeRange(chi.URLParam(r, "appId"), startTime, endTime)

	var apiStats dto.StatisticsPlot
	var peakConnectionStats dto.StatisticsPlot
	var websocketMessageStats dto.StatisticsPlot

	timestamps := stats.Timestamps()
	for i := len(timestamps) - 1; i >= 0; i-- {
		timestamp := timestamps[i]
		stat := stats.Get(timestamp)
		if stat == nil {
			continue
		}

		if len(apiStats.X) == 0 {
			apiStats.X = make([]int64, 0)
			apiStats.Y = make([]int, 0)
		}

		if len(peakConnectionStats.X) == 0 {
			peakConnectionStats.X = make([]int64, 0)
			peakConnectionStats.Y = make([]int, 0)
		}

		if len(websocketMessageStats.X) == 0 {
			websocketMessageStats.X = make([]int64, 0)
			websocketMessageStats.Y = make([]int, 0)
		}

		apiStats.X = append(apiStats.X, timestamp)
		peakConnectionStats.X = append(peakConnectionStats.X, timestamp)
		websocketMessageStats.X = append(websocketMessageStats.X, timestamp)

		apiStats.Y = append(apiStats.Y, stat.ApiMessages())
		peakConnectionStats.Y = append(peakConnectionStats.Y, stat.PeakConnections())
		websocketMessageStats.Y = append(websocketMessageStats.Y, stat.WebsocketMessages())
	}

	resp := dto.StatisticsGraph{
		ApiStats:              apiStats,
		PeakConnectionStats:   peakConnectionStats,
		WebsocketMessageStats: websocketMessageStats,
	}

	rendering.RenderSuccessWithData(w, "success", http.StatusOK, resp)
}
