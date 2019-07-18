package server

import (
	"github.com/jacoduplessis/golf"
	"log"
	"net/http"
	"time"
)

func intervalUpdate(c http.Client, tours []golf.Tour) {

	tickerLeaderboards := time.NewTicker(1 * time.Minute)
	go func() {
		for range tickerLeaderboards.C {
			updateLeaderboards(c, tours)
		}
	}()

	tickerTournaments := time.NewTicker(6 * time.Hour)
	go func() {
		for range tickerTournaments.C {
			updateTournaments(c, tours)
		}
	}()

	tickerTwitterClient := time.NewTicker(3 * time.Hour)
	go func() {
		for range tickerTwitterClient.C {
			setupTwitterClient(c)
		}
	}()

}

func updateTournaments(c http.Client, tours []golf.Tour) {

	errs := make(chan error, len(tours))

	for _, tour := range tours {
		go func(t golf.Tour) {
			errs <- t.UpdateTID(c)
		}(tour)
	}

	for range tours {
		err := <-errs
		if err != nil {
			log.Printf("Error updating TID: %s", err)
		}
	}
}

func updateLeaderboards(c http.Client, tours []golf.Tour) {

	results := make(chan *golf.Leaderboard, len(tours))

	for _, tour := range tours {

		go func(t golf.Tour) {
			r, err := t.Request()
			if err != nil {
				log.Printf("Error getting update request: %s\n", err)
				results <- nil
				return
			}
			resp, err := c.Do(r)
			if err != nil {
				log.Printf("Error performing update request: %s\n", err)
				results <- nil
				return
			}
			defer resp.Body.Close()
			lb, err := t.Parse(resp.Body)
			if err != nil {
				log.Printf("Error parsing leaderboard response: %s\n", err)
				results <- nil
				return
			}
			results <- lb

		}(tour)

	}

	for _, tour := range tours {
		lb := <-results
		if lb != nil {
			tour.SetLeaderboard(lb)
			tour.SetLastUpdated(time.Now())
		}
	}

}
