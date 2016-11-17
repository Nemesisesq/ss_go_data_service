package timers

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/nemesisesq/ss_data_service/streamsavvy"
)
func GuideboxEpisodeTimer () {

	quit := make(chan struct{})
	//ticker := time.NewTicker(24 * time.Hour)
	ticker := time.NewTicker(4 * time.Minute)
	go func(ticker *time.Ticker, quit chan struct{}) {
		log.WithFields(log.Fields{
			"timer": "ticker",
			"chan": "quit",
		}).Info("Launching Guidebox Episode Timer")
		for {
			select {
			case <-ticker.C:
				log.Println("ticker fired")
				streamsavvy.RefreshEpisodes()
			case <-quit:
				log.Println("ticker Stoping")
				ticker.Stop()
				return
			}
		}

		log.Println("Cleaning up!!")
	}(ticker, quit)

}
