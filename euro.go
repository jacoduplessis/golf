package main

import (
	"io"
	"net/http"
	"time"
)

// http://www.europeantour.com/library/MyEt/TourType=1/leaderboardJSON.js
// http://www.europeantour.com/data/tournament/2018050/leaderboard/languagecode/eng/isproam/0/feed/

type Euro struct {
	LastUpdated time.Time
	Leaderboard *Leaderboard
}

func (euro *Euro) Request() (*http.Request, error) {
	panic("implement me")
}

func (euro *Euro) Parse(io.Reader) (*Leaderboard, error) {
	panic("implement me")
}

func (euro *Euro) SetLeaderboard(lb *Leaderboard) {
	euro.Leaderboard = lb
}

func (euro *Euro) GetLeaderboard() *Leaderboard {
	return euro.Leaderboard
}

func (euro *Euro) GetLastUpdated() time.Time {
	return euro.LastUpdated
}

func (euro *Euro) SetLastUpdated(t time.Time) {
	euro.LastUpdated = t
}

func (euro *Euro) String() string {
	return "European Tour"
}

type EuroLeaderboard struct {
}
