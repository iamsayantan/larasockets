package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"github.com/iamsayantan/larasockets/statistics"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 6) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 100
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Connection encapsulates each incoming connection to our server.
type Connection struct {
	// id is the unique identifier for this connection.
	id string
	// app represents the application to which connection was made.
	app       *larasockets.Application
	collector statistics.StatsCollector

	logger        *zap.Logger
	websocketConn *websocket.Conn
	hub           *Hub

	// channel for outbound messages
	sendCh chan []byte
	// closeCh will be closed when the connection closes
	closeCh chan bool
}

// NewConnection generates a new Connection instance from the raw websocket connection.
func NewConnection(hub *Hub, app *larasockets.Application, conn *websocket.Conn, collector statistics.StatsCollector, logger *zap.Logger) larasockets.Connection {
	connId := generateIdForConnection()
	newConn := &Connection{
		id:            connId,
		app:           app,
		hub:           hub,
		collector:     collector,
		websocketConn: conn,
		logger:        logger.With(zap.String("application_id", app.Id()), zap.String("connection_id", connId)),
		sendCh:        make(chan []byte),
		closeCh:       make(chan bool),
	}

	go newConn.Receive()
	go newConn.writePump()

	newConn.collector.HandleConnection(app.Id())
	return newConn
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) Send(data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("error marshaling data to json",
			zap.String("error", err.Error()),
		)
		return
	}

	// c.collector.HandleWebsocketMessage(c.App().Id())
	c.sendCh <- message
}

func (c *Connection) Receive() {
	defer func() {
		c.Close()
	}()

	c.websocketConn.SetReadLimit(maxMessageSize)
	_ = c.websocketConn.SetReadDeadline(time.Now().Add(pongWait))

	c.websocketConn.SetPongHandler(func(string) error {
		_ = c.websocketConn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var pusherMessagePayload messages.PusherIncomingMessagePayload
		_, message, err := c.websocketConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("unexpected close error on websocket connection",
					zap.String("error", err.Error()),
				)
				return
			}

			c.logger.Error("error reading message from websocket connection",
				zap.String("error", err.Error()),
			)
			return
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = json.Unmarshal(message, &pusherMessagePayload)
		if err != nil {
			c.logger.Error("error unmarshalling received websocket message",
				zap.String("error", err.Error()),
			)
		}

		c.logger.Info("received message from websocket connection",
			zap.String("event", pusherMessagePayload.Event),
		)

		pusherMessage := messages.NewPusherMessage(c, c.hub.channelManger, pusherMessagePayload)
		pusherMessage.Respond()
	}
}

// writePump writes messages to the websocket connection.

// A goroutine running writePump is started for each connection. The application ensures
// there is at most one writer to a connection by executing all writes from this goroutine.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCh:
			_ = c.websocketConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// if we fail to read message from the send channel of the client, that means the hub
				// closed the connection. so we just inform the client that the connection has been closed.
				_ = c.websocketConn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			w, err := c.websocketConn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error("error acquiring next writer",
					zap.String("error", err.Error()),
				)
				continue
			}

			c.logger.Info("sending message",
				zap.String("message_payload", string(message)),
			)

			_, _ = w.Write(message)
			if err = w.Close(); err != nil {
				c.logger.Error("error closing the writer",
					zap.String("error", err.Error()),
				)
				continue
			}
		case <-ticker.C:
			_ = c.websocketConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.websocketConn.WriteMessage(websocket.PingMessage, []byte("PING")); err != nil {
				c.logger.Error("error writing ping message",
					zap.String("error", err.Error()),
				)
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

func (c *Connection) Close() {
	c.hub.RemoveConnection(c)
	c.closeCh <- true
	_ = c.websocketConn.Close()
	c.collector.HandleDisconnection(c.App().Id())
}

func (c *Connection) App() *larasockets.Application {
	return c.app
}

func generateIdForConnection() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(100000)) + "." + strconv.Itoa(rand.Intn(1000000000))
}
