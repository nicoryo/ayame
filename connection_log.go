package ayame

import (
	"slices"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func (c *connection) signalingLog(message message, rawMessage []byte) {
	// signaling の type を指定してフィルターできる
	// signaling_log_filters = register,offer,answer
	if !slices.Contains(c.config.SignalingLogFilters, message.Type) {
		return
	}

	if message.Type != "pong" {
		if c.config.Debug {
			c.signalingLogger.Debug().
				Str("roomId", c.roomID).
				Str("clientId", c.clientID).
				Str("connectionId", c.ID).
				Str("type", message.Type).
				Str("rawMessage", string(rawMessage)).
				Send()
		} else {
			c.signalingLogger.Info().
				Str("roomId", c.roomID).
				Str("clientId", c.clientID).
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
