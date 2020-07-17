package zlogging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// zerolog severity level
// panic (zerolog.PanicLevel, 5)
// fatal (zerolog.FatalLevel, 4)
// error (zerolog.ErrorLevel, 3)
// warn (zerolog.WarnLevel, 2)
// info (zerolog.InfoLevel, 1)
// debug (zerolog.DebugLevel, 0)
// trace (zerolog.TraceLevel, -1)

// Config used to configurate logger

var levelMapping = map[string]int{
	"trace": -1,
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
	"fatal": 4,
	"panic": 5,
}

type Config struct {
	Local bool   `yaml:"local"`
	Level string `yaml:"level"`
	AppID string `yaml:"app_id"`
	Env   string `yaml:"env"`
}

func (c Config) GetLevel() zerolog.Level {
	return zerolog.Level(levelMapping[c.Level])
}

func consoleLogger() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, FormatLevel: LevelFormater, TimeFormat: time.RFC3339}
	return log.Output(writer)
}

func SetupLogger(config Config) {
	var logger zerolog.Logger

	zerolog.DisableSampling(true)
	zerolog.ErrorStackMarshaler = MarshalStack
	hostname, _ := os.Hostname()

	zerolog.SetGlobalLevel(config.GetLevel())

	if config.Local {
		logger = consoleLogger()
	} else {
		logger = zerolog.New(os.Stdout)
	}

	logger = SetupHooks(logger)

	log.Logger = logger.With().Fields(
		map[string]interface{}{
			"app_id": config.AppID,
			"env":    config.Env,
			"host":   hostname,
		},
	).Timestamp().Logger()
}
