package server

import (
	"github.com/jacoduplessis/golf"
	"github.com/jacoduplessis/twitterparse"
	"log"
	"net/http"
	"sort"
)

func news(w http.ResponseWriter, r *http.Request) *AppError {

	results := make(chan []*twitterparse.Tweet, len(tours))
	var tweets []*twitterparse.Tweet

	for _, tour := range tours {

		go func(t golf.Tour) {
			twitterID := t.TwitterID()
			if twitterID == "" {
				results <- nil
				return
			}

			tweets, err := twitterClient.GetProfileTweets(twitterID)
			if err != nil {
				log.Printf("error fetching tweets for %s %s", t, err)
			}
			results <- tweets
		}(tour)

	}

	for range tours {
		tweets = append(tweets, <-results...)
	}

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].Timestamp > tweets[j].Timestamp
	})

	newsTemplate.ExecuteTemplate(w, "", tweets)
	return nil
}

func index(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := struct {
		Leaderboards []*golf.Leaderboard
	}{}

	for _, tour := range tours {
		if len(tour.Leaderboard().Players) > 0 {
			ctx.Leaderboards = append(ctx.Leaderboards, tour.Leaderboard())
		}
	}

	sort.Slice(ctx.Leaderboards, func(i, j int) bool {
		return ctx.Leaderboards[i].TourIndex < ctx.Leaderboards[j].TourIndex
	})

	if r.URL.Query().Get("format") == "json" {
		return renderJSON(w, ctx)
	}
	return renderTemplate(w, ctx)

}
