package ayame

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// megabytes
	logRotateMaxSize    = 10
	logRotateMaxBackups = 5
	//days
	logRotateMaxAge = 30
)

func InitLogger(config *Config) error {
	if f, err := os.Stat(config.LogDir); os.IsNotExist(err) || !f.IsDir() {
		return err
	}

	logPath := fmt.Sprintf("%s/%s", config.LogDir, config.LogName)

	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    logRotateMaxSize,
		MaxBackups: logRotateMaxBackups,
		MaxAge:     logRotateMaxAge,
		Compress:   true,
	}

	// https://github.com/rs/zerolog/issues/77
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logLevel, err := parseLevel(*config)
	if err != nil {
		return err
	}

	// デバッグが有効な時はコンソールにもだす
	if config.Debug {
		var writers io.Writer
		stdout := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: "2006-01-02 15:04:05.000000Z"}
		prettyFormat(&stdout)
		writers = io.MultiWriter(stdout, writer)
		log.Logger = zerolog.New(writers).With().Caller().Timestamp().Logger().Level(logLevel)
	} else {
		log.Logger = zerolog.New(writer).With().Timestamp().Logger().Level(logLevel)
	}

	return nil
}

// 現時点での prettyFormat
// 2023-04-17 12:51:56.333485Z [INFO] config.go:102 > CONF | debug=true
func prettyFormat(w *zerolog.ConsoleWriter) {
	const Reset = "\x1b[0m"

	w.FormatLevel = func(i interface{}) string {
		var color, level string
		// TODO: 各色を定数に置き換える
		// TODO: 他の logLevel が必要な場合は追加する
		switch i.(string) {
		case "info":
			color = "\x1b[32m"
		case "error":
			color = "\x1b[31m"
		case "warn":
			color = "\x1b[33m"
		case "debug":
			color = "\x1b[34m"
		default:
			color = "\x1b[37m"
		}

		level = strings.ToUpper(i.(string))
		return fmt.Sprintf("%s[%s]%s", color, level, Reset)
	}
	w.FormatCaller = func(i interface{}) string {
		return fmt.Sprintf("[%s]", filepath.Base(i.(string)))
	}
	// TODO: Caller をファイル名と行番号だけの表示で出力する
	//       以下のようなフォーマットにしたい
	//       2023-04-17 12:50:09.334758Z [INFO] [config.go:102] CONF | debug=true
	// TODO: name=value が無い場合に | を消す方法がわからなかった
	w.FormatMessage = func(i interface{}) string {
		if i == nil || i == "" {
			return ""
		}
		return fmt.Sprintf("%s |", i)
	}
	w.FormatFieldName = func(i interface{}) string {
		const Cyan = "\x1b[36m"
		return fmt.Sprintf("%s%s=%s", Cyan, i, Reset)
	}
	// TODO: カンマ区切りを同実現するかわからなかった
	w.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
}

func parseLevel(config Config) (zerolog.Level, error) {
	// debug: true の場合の log_level は debug で固定
	if config.Debug {
		return zerolog.DebugLevel, nil
	}

	// 空文字列は NoLevel 扱いで ParseLevel でエラーにならないため事前に確認する
	if config.LogLevel == "" {
		return zerolog.NoLevel, errConfigInvalidLogLevel
	}

	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		// err は継続するように読めるのでここで捨てる
		return logLevel, errConfigInvalidLogLevel
	}

	return logLevel, nil
}
