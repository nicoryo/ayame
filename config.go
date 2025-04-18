package ayame

import (
	_ "embed"
	"net/url"

	zlog "github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
)

//go:embed VERSION
var Version string

const (
	defaultLogDir              = "."
	defaultLogName             = "ayame.jsonl"
	defaultLogRotateMaxSize    = 200
	defaultLogRotateMaxBackups = 7
	defaultLogRotateMaxAge     = 30
	defaultLogRotateCompress   = false

	defaultSignalingLogName = "signaling.jsonl"

	defaultWebSocketReadTimeoutSec  = 90
	defaultWebSocketPongTimeoutSec  = 60
	defaultWebSocketPingIntervalSec = 5

	defaultWebhookLogName        = "webhook.jsonl"
	defaultWebhookRequestTimeout = 5

	defaultListenPrometheusIPv4Address = "0.0.0.0"
	defaultListenPrometheusPortNumber  = 4000
)

type Config struct {
	Debug bool `ini:"debug"`

	LogDir              string `ini:"log_dir"`
	LogName             string `ini:"log_name"`
	LogStdout           bool   `ini:"log_stdout"`
	LogRotateMaxSize    int    `ini:"log_rotate_max_size"`
	LogRotateMaxBackups int    `ini:"log_rotate_max_backups"`
	LogRotateMaxAge     int    `ini:"log_rotate_max_age"`
	LogRotateCompress   bool   `ini:"log_rotate_compress"`

	SignalingLogName string `ini:"signaling_log_name"`

	DebugConsoleLog     bool `ini:"debug_console_log"`
	DebugConsoleLogJSON bool `ini:"debug_console_log_json"`

	TypeMessage bool `ini:"type_message"`

	ListenIPv4Address string `ini:"listen_ipv4_address"`
	ListenPortNumber  int32  `ini:"listen_port_number"`

	// socket の待ち受け時間
	WebSocketReadTimeoutSec int32 `ini:"websocket_read_timeout_sec"`
	// pong が送られてこないためタイムアウトにするまでの時間
	WebSocketPongTimeoutSec int32 `ini:"websocket_pong_timeout_sec"`
	// ping 送信の時間間隔
	WebSocketPingIntervalSec int32 `ini:"websocket_ping_interval_sec"`

	AuthnWebhookURL      string `ini:"authn_webhook_url"`
	DisconnectWebhookURL string `ini:"disconnect_webhook_url"`

	WebhookLogName           string `ini:"webhook_log_name"`
	WebhookRequestTimeoutSec int32  `ini:"webhook_request_timeout_sec"`

	ListenPrometheusIPv4Address string `ini:"listen_prometheus_ipv4_address"`
	ListenPrometheusPortNumber  int32  `ini:"listen_prometheus_port_number"`
}

func NewConfig(configFilePath string) (*Config, error) {
	config := new(Config)

	iniConfig, err := ini.InsensitiveLoad(configFilePath)
	if err != nil {
		return nil, err
	}

	if err := iniConfig.StrictMapTo(config); err != nil {
		return nil, err
	}

	if config.AuthnWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
			return nil, err
		}
	}

	if config.DisconnectWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.DisconnectWebhookURL); err != nil {
			return nil, err
		}
	}

	setDefaultsConfig(config)

	return config, nil
}

func setDefaultsConfig(config *Config) {
	if config.LogDir == "" {
		config.LogDir = defaultLogDir
	}

	if config.LogName == "" {
		config.LogDir = defaultLogName
	}

	if config.SignalingLogName == "" {
		config.SignalingLogName = defaultSignalingLogName
	}

	if config.WebSocketReadTimeoutSec == 0 {
		config.WebSocketReadTimeoutSec = defaultWebSocketReadTimeoutSec
	}

	if config.WebSocketPongTimeoutSec == 0 {
		config.WebSocketPongTimeoutSec = defaultWebSocketPongTimeoutSec
	}

	if config.WebSocketPingIntervalSec == 0 {
		config.WebSocketPingIntervalSec = defaultWebSocketPingIntervalSec
	}

	if config.WebhookLogName == "" {
		config.WebhookLogName = defaultWebhookLogName
	}

	if config.WebhookRequestTimeoutSec == 0 {
		config.WebhookRequestTimeoutSec = defaultWebhookRequestTimeout
	}

	if config.ListenPrometheusIPv4Address == "" {
		config.ListenPrometheusIPv4Address = defaultListenPrometheusIPv4Address
	}

	if config.ListenPrometheusPortNumber == 0 {
		config.ListenPrometheusPortNumber = defaultListenPrometheusPortNumber
	}

}

func (c *Config) PrintConfig() {
	zlog.Info().Bool("debug", c.Debug).Msg("AyameConf")

	zlog.Info().Str("log_dir", c.LogDir).Msg("AyameConf")
	zlog.Info().Str("log_name", c.LogName).Msg("AyameConf")
	zlog.Info().Bool("log_stdout", c.LogStdout).Msg("AyameConf")

	zlog.Info().Int("log_rotate_max_size", c.LogRotateMaxSize).Msg("AyameConf")
	zlog.Info().Int("log_rotate_max_backups", c.LogRotateMaxBackups).Msg("AyameConf")
	zlog.Info().Int("log_rotate_max_age", c.LogRotateMaxAge).Msg("AyameConf")
	zlog.Info().Bool("log_rotate_compress", c.LogRotateCompress).Msg("AyameConf")

	zlog.Info().Str("signaling_log_name", c.SignalingLogName).Msg("AyameConf")

	zlog.Info().Bool("debug_console_log", c.DebugConsoleLog).Msg("AyameConf")
	zlog.Info().Bool("debug_console_log_json", c.DebugConsoleLogJSON).Msg("AyameConf")

	zlog.Info().Str("listen_ipv4_address", c.ListenIPv4Address).Msg("AyameConf")
	zlog.Info().Int32("listen_port_number", c.ListenPortNumber).Msg("AyameConf")

	zlog.Info().Str("authn_webhook_url", c.AuthnWebhookURL).Msg("AyameConf")
	zlog.Info().Str("disconnect_webhook_url", c.DisconnectWebhookURL).Msg("AyameConf")

	zlog.Info().Str("webhook_log_name", c.WebhookLogName).Msg("AyameConf")
	zlog.Info().Int32("webhook_request_timeout_sec", c.WebhookRequestTimeoutSec).Msg("AyameConf")

	zlog.Info().Str("prometheus_ipv4_address", c.ListenPrometheusIPv4Address).Msg("AyameConf")
	zlog.Info().Int32("prometheus_port", c.ListenPrometheusPortNumber).Msg("AyameConf")
}
