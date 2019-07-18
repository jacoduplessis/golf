package server

import (
	"github.com/jacoduplessis/twitterparse"
	"log"
	"net/http"
)

func setupTwitterClient(c http.Client) {
	tc, err := twitterparse.NewClientWithHTTPClient(c)
	if err != nil {
		log.Fatalf("Error setting up twitter client: %v\n", err)
	}
	twitterClient = tc
}
