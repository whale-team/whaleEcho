package zlogging

import (
	"path"
	"runtime"

	"github.com/rs/zerolog"
)

// CallerHook add caller field into log message
type CallerHook struct{}

// Run run the hook
func (h CallerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if _, file, line, ok := runtime.Caller(3); ok {
		e.Str("file", path.Base(file)).Int("line", line)
	}
}

// SetupHooks setup all hooks
func SetupHooks(logger zerolog.Logger) zerolog.Logger {
	return logger.Hook(CallerHook{}).With().Logger()
}
