package ayame

import (
	"github.com/rs/zerolog"
)

func InitWebhookLogger(config *Config) (*zerolog.Logger, error) {
	return InitLogger(config, config.WebhookLogName, "webhook")
}
