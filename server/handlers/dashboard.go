package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/server/handlers/dto"
	"github.com/iamsayantan/larasockets/server/handlers/middlewares"
	"github.com/iamsayantan/larasockets/server/rendering"
	"github.com/pusher/pusher-http-go/v5"
	"io/ioutil"
	"net/http"
	"time"
)

func NewDashboardHandler(cm larasockets.ChannelManager) *DashboardHandler {
	return &DashboardHandler{appManager: cm.AppManager(), channelManager: cm}
}

type DashboardHandler struct {
	appManager     larasockets.ApplicationManager
	channelManager larasockets.ChannelManager
}

func (h *DashboardHandler) AllApps(w http.ResponseWriter, r *http.Request) {
	apps := h.appManager.All()
	appResponses := make([]dto.ApplicationResponse, 0)

	for _, app := range apps {
		appResponse := dto.ApplicationResponse{
			AppId:   app.Id(),
			AppName: app.Name(),
			AppKey:  app.Key(),
		}

		appResponses = append(appResponses, appResponse)
	}

	rendering.RenderSuccessWithData(w, "success", http.StatusOK, appResponses)
}

func (h *DashboardHandler) AuthorizeConnectionRequest(w http.ResponseWriter, r *http.Request) {
	var connectionRequest dto.ConnectionRequest

	err := json.NewDecoder(r.Body).Decode(&connectionRequest)
	if err != nil {
		rendering.RenderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := h.appManager.FindById(connectionRequest.AppId)
	if app == nil {
		rendering.RenderError(w, "invalid app selected", http.StatusBadRequest)
		return
	}

	if app.Secret() != connectionRequest.AppSecret {
		rendering.RenderError(w, "invalid app secret", http.StatusBadRequest)
		return
	}

	jwtExpirationTime := time.Now().Add(time.Hour * 24)
	jwtClaims := dto.ApplicationAuthorizationClaims{
		AppId: app.Id(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtExpirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessToken, err := jwtToken.SignedString([]byte(app.Secret()))
	if err != nil {
		rendering.RenderError(w, "error generating access token", http.StatusInternalServerError)
		return
	}

	resp := dto.ConnectionAuthorizationResponse{
		AppId:       app.Id(),
		ApiKey:      app.Key(),
		AccessToken: accessToken,
	}

	rendering.RenderSuccessWithData(w, "success", http.StatusOK, resp)
}

func (h *DashboardHandler) AuthorizeChannelRequest(w http.ResponseWriter, r *http.Request) {
	appId := middlewares.GetAuthenticatedAppIdFromContext(r.Context())
	app := h.appManager.FindById(appId) // the appId will be already validated in the middleware, skipping nil check.

	pusherClient := pusher.Client{
		AppID:  app.Id(),
		Key:    app.Key(),
		Secret: app.Secret(),
	}

	params, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rendering.RenderError(w, "error reading request body", http.StatusInternalServerError)
		return
	}

	response, err := pusherClient.AuthenticatePrivateChannel(params)
	if err != nil {
		rendering.RenderError(w, "error authorizing channel", http.StatusBadRequest)
		return
	}

	_, _ = w.Write(response)
}

func (h *DashboardHandler) TriggerEvent(w http.ResponseWriter, r *http.Request) {
	appId := middlewares.GetAuthenticatedAppIdFromContext(r.Context())

	var triggerEventRequest dto.DashboardEventTriggerRequest
	err := json.NewDecoder(r.Body).Decode(&triggerEventRequest)
	if err != nil {
		rendering.RenderError(w, fmt.Sprintf("invalid request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	channel := h.channelManager.FindChannel(appId, triggerEventRequest.Channel)
	if channel == nil {
		return
	}

	messagePayload := messages.PusherEventPayload{
		Channel: triggerEventRequest.Channel,
		Event:   triggerEventRequest.Event,
		Data:    triggerEventRequest.Data,
	}

	channel.Broadcast(messagePayload)
}
