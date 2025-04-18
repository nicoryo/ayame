package ayame

import (
	"github.com/rs/zerolog"
)

func InitSignalingLogger(config *Config) (*zerolog.Logger, error) {
	return InitLogger(config, config.SignalingLogName, "signaling")
}
