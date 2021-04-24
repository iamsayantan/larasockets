package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/events"
	"github.com/iamsayantan/larasockets/server/handlers"
	"github.com/iamsayantan/larasockets/statistics"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	logger *zap.Logger
	router chi.Router
	hub    *Hub

	collector      statistics.StatsCollector
	channelManager larasockets.ChannelManager
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request) {
	appKey := chi.URLParam(r, "appKey")
	if appKey == "" {
		s.logger.Error("no appKey found on url path")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	app := s.channelManager.AppManager().FindByKey(appKey)
	if app == nil {
		s.logger.Error("invalid appKey. no app found with the given appKey", zap.String("appKey", appKey))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("error upgrading to websocket connection", zap.String("error", err.Error()))
		return
	}

	wsConn := NewConnection(s.hub, app, conn, s.collector, s.logger)
	s.hub.register <- wsConn

	s.logger.Info("received new websocket connection", zap.String("connection_id", wsConn.Id()), zap.String("application_id", app.Id()))

	connResp := events.NewConnectionEstablished(wsConn.Id(), 120)
	wsConn.Send(connResp)
}

func NewServer(logger *zap.Logger, cm larasockets.ChannelManager, collector statistics.StatsCollector) *Server {
	server := &Server{}

	server.channelManager = cm
	server.logger = logger
	server.collector = collector
	server.hub = NewHub(logger, cm)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	triggerHandler := handlers.NewTriggerEventHandler(server.channelManager, server.collector, server.logger)
	statsHandler := handlers.NewStatsHandler(collector)

	r.Get("/app/{appKey}", server.ServeWS)
	r.Post("/apps/{appId}/events", triggerHandler.HandleEvents)
	r.Get("/apps/{appId}/stats", statsHandler.GetStatsForApp)

	server.router = r

	return server
}
