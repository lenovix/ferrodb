package engine

import (
	"fmt"
	"runtime"
	"time"
)

func (e *Engine) Info() string {
	uptime := time.Since(e.startTime).Seconds()

	return fmt.Sprintf(
		"FerroDB v0.3.0\n"+
			"uptime_seconds: %.0f\n"+
			"keys: %d\n"+
			"goroutines: %d\n"+
			"go_version: %s",
		uptime,
		e.store.Size(),
		runtime.NumGoroutine(),
		runtime.Version(),
	)
}
