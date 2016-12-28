package timers

import (
	log "github.com/Sirupsen/logrus"
	gnote "github.com/nemesisesq/ss_data_service/gracenote"
	"time"
)

func GraceNoteListingTimer() {

	quit := make(chan struct{})
	ticker := time.NewTicker(15 * time.Minute)
	go func(ticker *time.Ticker, quit chan struct{}) {
		log.WithFields(log.Fields{
			"timer": "ticker",
			"chan":  "quit",
		}).Info("Launching Gracenote Listing Timer")
		for {
			select {
			case <-ticker.C:
				log.Println("ticker fired")
				gnote.RefreshListings()
			case <-quit:
				log.Println("ticker Stoping")
				ticker.Stop()
				return
			}
		}

		log.Println("Cleaning up!!")
	}(ticker, quit)
}
