package ayame

import (
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func (c *connection) signalingLog(message message, rawMessage []byte) {
	if message.Type != "pong" {
		if c.config.Debug {
			c.signalingLogger.Log().
				Str("roomId", c.roomID).
				Str("clientID", c.clientID).
				Str("connectionId", c.ID).
				Str("type", message.Type).
				Msg(string(rawMessage))
		} else {
			c.signalingLogger.Log().
				Str("roomId", c.roomID).
				Str("clientID", c.clientID).
				Str("connectionId", c.ID).
				Str("type", message.Type).
				Send()
		}
	}
}

func (c *connection) errLog() *zerolog.Event {
	return zlog.Error().
		Str("roomId", c.roomID).
		Str("clientID", c.clientID).
		Str("connectionId", c.ID)
}

func (c *connection) debugLog() *zerolog.Event {
	return zlog.Debug().
		Str("roomId", c.roomID).
		Str("clientID", c.clientID).
		Str("connectionId", c.ID)
}
