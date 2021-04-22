package channels

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/messages"
	"log"
	"strings"
)

type privateChannel struct {
	publicChannel
}

func (c *privateChannel) Subscribe(conn larasockets.Connection, payload interface{}) {
	subscriptionPayload, ok := payload.(messages.PusherSubscriptionPayload)
	if !ok {
		log.Printf("error converting the payload")
		return
	}

	err := c.verifySignature(conn, subscriptionPayload)
	if err != nil {
		log.Printf("error verifying signature: %s", err.Error())
		// see https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol#Error-Codes
		errMessage := messages.NewPusherErrorMessage(err.Error(), 4009)
		conn.Send(errMessage)

		return
	}

	if c.IsSubscribed(conn) {
		log.Printf("connection already subscribed")
		return
	}

	c.connections[conn.Id()] = conn

	resp := struct {
		Event   string      `json:"event"`
		Channel string      `json:"channel"`
		Data    interface{} `json:"data"`
	}{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: c.Name(),
		Data:    "{}",
	}
	conn.Send(resp)
}

func (c *privateChannel) verifySignature(conn larasockets.Connection, payload messages.PusherSubscriptionPayload) error {
	// The signature is generated in the following format "<socket-id>:<channel-name>" and then its signed using
	// the secret key. So we will generate the signature ourselves and verify if it matches with us.
	var signature strings.Builder
	signature.WriteString(conn.Id())
	signature.WriteString(":")
	signature.WriteString(c.Name())

	if payload.ChannelData != "" {
		signature.WriteString(":")
		signature.WriteString(payload.ChannelData)
	}

	// signature format "<pusher key>:<signature>". So we need to validate the signature part.
	s := strings.SplitAfter(payload.Auth, ":")
	incomingSignature, _ := hex.DecodeString(strings.Join(s[1:], ""))

	signedSignature := signature.String()
	h := hmac.New(sha256.New, []byte(conn.App().Secret()))
	h.Write([]byte(signedSignature))

	if valid := hmac.Equal(incomingSignature, h.Sum(nil)); !valid {
		return errors.New("invalid auth signature")
	}

	return nil
}

func newPrivateChannel(name string) larasockets.Channel {
	return &privateChannel{publicChannel{
		name:        name,
		connections: make(map[string]larasockets.Connection, 0),
	}}
}
