package golf

import (
	"io"
	"net/http"
	"time"
)

type Player struct {
	Name            string
	Country         string
	CurrentPosition string
	StartPosition   string
	After           int // holes completed
	Hole            int // current hole
	Today           int // par
	Total           int // par
	TotalStrokes    int // stroke
	Rounds          []int
}

type Leaderboard struct {
	Tour       string
	TourIndex  int
	Tournament string
	Round      int
	Location   string
	Course     string
	Date       string
	Updated    string
	Players    []*Player
}

type Tour interface {
	Request() (*http.Request, error)
	Parse(io.Reader) (*Leaderboard, error)
	Leaderboard() *Leaderboard
	SetLeaderboard(*Leaderboard)
	LastUpdated() time.Time
	SetLastUpdated(time.Time)
	TID() string
	UpdateTID(http.Client) error
	Twitter() string
	TwitterID() string
	Index() int
}
