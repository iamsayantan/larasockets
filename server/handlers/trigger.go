package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/events"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/statistics"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sort"
	"strings"
)

// URL /apps/{appId}/events
type TriggerEventsHandler struct {
	channelManager larasockets.ChannelManager
	collector      statistics.StatsCollector

	logger *zap.Logger
}

type PusherServerEventPayload struct {
	Name     string   `json:"name"`
	Data     string   `json:"data"`
	Channel  string   `json:"channel"`
	Channels []string `json:"channels"`
	SocketId string   `json:"socket_id"`
}

func NewTriggerEventHandler(cm larasockets.ChannelManager, collector statistics.StatsCollector, logger *zap.Logger) *TriggerEventsHandler {
	return &TriggerEventsHandler{channelManager: cm, logger: logger.With(zap.String("handler", "TriggerEventsHandler")), collector: collector}
}

func (h *TriggerEventsHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("received request on HandleEvents")

	appId := chi.URLParam(r, "appId")
	var bodyParams PusherServerEventPayload
	err := h.verifySignature(r)
	if err != nil {
		h.logger.Error("error verifying authentication signature", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
	}

	err = json.NewDecoder(r.Body).Decode(&bodyParams)
	if err != nil {
		h.logger.Error("error decoding json request.", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	h.collector.HandleApiMessage(appId)

	for _, channelName := range bodyParams.Channels {
		channel := h.channelManager.FindChannel(appId, channelName)
		if channel == nil {
			h.logger.Info("channel not found", zap.String("channel_name", channelName), zap.String("application_id", appId))
			log.Printf("channel %s not found for app %s", channelName, appId)
			continue
		}

		payload := messages.PusherEventPayload{
			Event:   bodyParams.Name,
			Channel: channelName,
			Data:    bodyParams.Data,
		}

		events.LogEvent(h.channelManager, events.ApiMessage, events.DashboardLogDetails{
			AppId:        appId,
			ChannelName:  channelName,
			EventName:    bodyParams.Name,
			ConnectionId: "",
			EventPayload: bodyParams.Data,
		})

		if bodyParams.SocketId == "" {
			channel.Broadcast(payload)
		} else {
			channel.BroadcastExcept(payload, bodyParams.SocketId)
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}

// verifySignature validates the auth signature from the incoming request.
// See https://pusher.com/docs/channels/library_auth_reference/rest-api#Authentication for more implementation
// details.
func (h *TriggerEventsHandler) verifySignature(r *http.Request) error {
	// Whether the appId is valid will be checked in the middleware layer, so at this
	// point it is safe to assume that a valid app exists with the given appId.
	app := h.channelManager.AppManager().FindById(chi.URLParam(r, "appId"))

	queryParams := r.URL.Query()
	queryParamsKeys := make([]string, 0)
	authSignature := queryParams.Get("auth_signature")

	for key := range queryParams {
		if key == "auth_signature" {
			continue
		}

		queryParamsKeys = append(queryParamsKeys, key)
	}
	sort.Strings(queryParamsKeys)

	var signatureString strings.Builder
	sortedQueryParams := make([]string, 0)

	signatureString.WriteString(r.Method)
	signatureString.WriteString("\n")
	signatureString.WriteString(r.URL.Path)
	signatureString.WriteString("\n")

	for _, key := range queryParamsKeys {
		var str strings.Builder
		str.WriteString(key)
		str.WriteString("=")
		str.WriteString(queryParams.Get(key))

		sortedQueryParams = append(sortedQueryParams, str.String())
	}

	signatureString.WriteString(strings.Join(sortedQueryParams, "&"))

	incomingSignature, err := hex.DecodeString(authSignature)
	if err != nil {
		return err
	}

	hashToSign := hmac.New(sha256.New, []byte(app.Secret()))
	hashToSign.Write([]byte(signatureString.String()))

	if valid := hmac.Equal(hashToSign.Sum(nil), incomingSignature); !valid {
		return errors.New("invalid auth signature")
	}

	return nil
}
