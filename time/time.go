package time

import (
	"net/http"
	"time"

	"anidiot.com/auto-trade/log"
	"go.uber.org/zap"
	// "log"
)

var timer <-chan time.Time
var closeSignal <-chan bool
var client http.Client

func InitWatcher(duration time.Duration, stop <-chan bool, handler func(time.Time)) {
	ticker := time.NewTicker(duration)
	client = http.Client{
		Timeout: time.Second * 2,
	}
	go func() {
		log.Logger.Info("Init time watcher")
		now := time.Now()
		handler(now)
		for {
			select {
			case <-stop:
				return
			case t := <-ticker.C:
				log.Logger.Info("", zap.Any("ticker fire", t))
				handler(t)
			}
		}
	}()
}
