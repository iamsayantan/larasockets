package handlers

import (
	"github.com/go-chi/chi"
	"github.com/iamsayantan/larasockets/server/handlers/dto"
	"github.com/iamsayantan/larasockets/server/rendering"
	"github.com/iamsayantan/larasockets/statistics"
	"net/http"
)

func NewStatsHandler(store statistics.StatsStorage) *StatsHandler {
	return &StatsHandler{statsStore: store}
}

type StatsHandler struct {
	statsStore statistics.StatsStorage
}

func (h *StatsHandler) GetStatForToday(w http.ResponseWriter, r *http.Request) {
	stat := h.statsStore.DailyStatForApp(chi.URLParam(r, "appId"))
	resp := dto.DailyStatSnapshot{
		PeakConnections:   stat.PeakConnections(),
		ApiMessages:       stat.ApiMessages(),
		WebsocketMessages: stat.WebsocketMessages(),
	}

	rendering.RenderSuccessWithData(w, "success", http.StatusOK, resp)
}
