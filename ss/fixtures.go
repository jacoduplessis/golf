package ss

import (
	"log"
	"os"
)

func GetMatchesFixture() []Match {

	r, err := os.Open("test/ss_results.html")
	if err != nil {
		log.Fatal(err)
	}
	m, err := ParseMatches(r)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func GetMatchFixture() Match {
	r, err := os.Open("test/ss_match.html")
	if err != nil {
		log.Fatal(err)
	}
	m, err := ParseMatch(r, "")
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func GetScorecardFixture() Scorecard {
	r, err := os.Open("test/ss_scorecard.json")
	if err != nil {
		log.Fatal(err)
	}
	sc, err := ParseScorecard(r)
	if err != nil {
		log.Fatal(err)
	}
	return sc
}
