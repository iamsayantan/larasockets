package handlers

import (
	"github.com/go-chi/chi"
	"github.com/iamsayantan/larasockets/server/rendering"
	"github.com/iamsayantan/larasockets/statistics"
	"net/http"
)

func NewStatsHandler(collector statistics.StatsCollector) *StatsHandler {
	return &StatsHandler{collector: collector}
}

type StatsHandler struct {
	collector statistics.StatsCollector
}

func (h *StatsHandler) GetStatsForApp(w http.ResponseWriter, r *http.Request) {
	stat := h.collector.GetAppStatistics(chi.URLParam(r, "appId"))
	rendering.RenderSuccessWithData(w, "success", http.StatusOK, stat.GetCurrentSnapshot())
}
