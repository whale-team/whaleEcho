package cmd

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
)

func cmdRecover() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg = msg + fmt.Sprintf("%s:%d\n", file, line)
		}
		log.Error().Stack().Err(r.(error)).Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
}
