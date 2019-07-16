package ss

import (
	"encoding/json"
	"errors"
	"fmt"
	q "github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
	"time"
)

type Int64Str int64

func (i Int64Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

func (i *Int64Str) UnmarshalJSON(b []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		var value int64
		if strings.ToLower(s) == "e" || s == "" {
			value = 0
		} else {
			value, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				value = 0
			}
		}

		*i = Int64Str(value)
		return nil
	}

	// Fallback to number
	return json.Unmarshal(b, (*int64)(i))
}

type Player struct {
	Name        string
	FirstName   string
	LastName    string
	Position    Int64Str
	Status      string
	Score       Int64Str
	Hole        Int64Str
	Round       Int64Str
	Today       string
	Rounds      []Int64Str
	Strokes     Int64Str
	ScorecardId string `json:"scorecard_id"`
	MatchId     string
}

func (p Player) StrRounds() string {
	s := ""
	for _, r := range p.Rounds {
		s = s + fmt.Sprintf("%d ", r)
	}
	return strings.TrimSpace(s)
}

type Match struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	TourName     string    `json:"tournament"`
	Location     string    `json:"location"`
	StartDate    time.Time `json:"starts_at"`
	EndDate      time.Time `json:"ends_at"`
	CurrentRound int       `json:"current_round"`
	PrizeMoney   string    `json:"prize_money"`
	Players      []*Player `json:"people"`
}

type Hole struct {
	Number  Int64Str
	Par     Int64Str
	Strokes Int64Str
}

func (h Hole) Result() string {

	if h.Strokes == 1 {
		return "ace"
	}

	if h.Par == h.Strokes+3 {
		return "albatross"
	}

	if h.Par == h.Strokes+2 {
		return "eagle"
	}

	if h.Par == h.Strokes+1 {
		return "birdie"
	}

	if h.Par == h.Strokes {
		return "par"
	}

	if h.Par == h.Strokes-1 {
		return "bogey"
	}

	if h.Strokes > h.Par+1 {
		return "disaster"
	}

	return ""
}

type Scorecard struct {
	Rounds []struct {
		Number  int
		Par     Int64Str
		Strokes Int64Str
		Holes   []Hole
	}
}

func ParseScorecard(r io.ReadCloser) (Scorecard, error) {

	var sc Scorecard
	err := json.NewDecoder(r).Decode(&sc)
	_ = r.Close()
	return sc, err
}

func ParseMatch(r io.ReadCloser, matchId string) (Match, error) {

	var match Match

	doc, err := q.NewDocumentFromReader(r)
	if err != nil {
		return match, err
	}
	_ = r.Close()

	content, ok := doc.Find("golf-match-details").Attr("match")
	if !ok {
		return match, errors.New("Parse error: missing attribute 'match'")
	}

	err = json.Unmarshal([]byte(content), &match)

	if match.ID == "" && matchId != "" {
		match.ID = matchId
	}

	for _, p := range match.Players {
		p.MatchId = match.ID
	}

	return match, err

}

func ParseMatches(r io.ReadCloser) ([]Match, error) {

	doc, err := q.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	_ = r.Close()

	el := doc.Find("golf-home-results")

	content, ok := el.Attr("matches")

	if !ok {
		return nil, errors.New("Could not parse response: missing attribute 'matches'")
	}

	matches := []Match{}

	err = json.Unmarshal([]byte(content), &matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
