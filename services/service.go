package services

import (
	"fmt"
	"os"
	"time"

	"github.com/littlehawk93/columba/config"
)

// BackgroundService a function that gets called repeatedly in the background by the background service scheduler
type BackgroundService func(cfg config.ApplicationConfiguration) error

// RunBackgroundService executes a background task repeatedly on a set interval
func RunBackgroundService(name string, service BackgroundService, cfg config.ApplicationConfiguration, interval time.Duration) {

	go func() {
		for {
			time.Sleep(interval)

			if err := service(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] Background Task %s: %s\n", time.Now().Format(time.RFC3339), name, err.Error())
			}
		}
	}()
}
