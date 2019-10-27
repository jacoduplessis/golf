package euro

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/cascadia"
	"github.com/jacoduplessis/golf"
	"github.com/jacoduplessis/simplejson"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// http://www.europeantour.com/library/MyEt/TourType=1/leaderboardJSON.js
// http://www.europeantour.com/data/tournament/2018050/leaderboard/languagecode/eng/isproam/0/feed/

// can also get tid here
// http://app.europeantour.com/mobile/

type Euro struct {
	lastUpdated time.Time
	leaderboard *golf.Leaderboard
	tid         string
}

var tidSelector = cascadia.MustCompile("mini-leaderboard")

func (euro *Euro) TID() string {
	return euro.tid
}

func (euro *Euro) Index() int {
	return 2
}

func extractTID(r io.Reader) (string, error) {

	var lData string
	root, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	d := tidSelector.MatchFirst(root)

	if d == nil {
		return "", fmt.Errorf("[euro] could not find html element")
	}

	for _, a := range d.Attr {
		if a.Key == ":event-data" {
			lData = a.Val
		}
	}
	if lData == "" {
		return "", errors.New("Euro TID is empty")
	}

	data, err := simplejson.NewJson([]byte(lData))
	if err != nil {
		return "", fmt.Errorf("error parsing JSON attribute: %s", err)
	}

	tidInt, err := data.Get("EventId").Int()
	return strconv.Itoa(tidInt), err

}

func (euro *Euro) UpdateTID(c http.Client) error {

	res, err := c.Get("http://www.europeantour.com")

	if err != nil {
		return fmt.Errorf("error fetching html from euro: %s", err)
	}

	defer res.Body.Close()

	tid, err := extractTID(res.Body)
	euro.tid = tid
	return err
}

func (euro *Euro) Request() (*http.Request, error) {

	u := fmt.Sprintf("https://www.europeantour.com/api/sportdata/Leaderboard/Strokeplay/%s", euro.TID())
	return http.NewRequest("GET", u, nil)
}

func (euro *Euro) Parse(r io.Reader) (*golf.Leaderboard, error) {
	var d EuroLeaderboard

	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	var players []*golf.Player

	for _, p := range d.Players {

		rounds := []int{}

		var totalStrokes int

		for _, r := range p.Rounds {
			rounds = append(rounds, r.Strokes)
			totalStrokes += r.Strokes
		}

		players = append(players, &golf.Player{
			Name:            FixEuroName(p.FirstName + " " + p.LastName),
			Country:         "",
			CurrentPosition: p.PositionDesc,
			StartPosition:   "",
			After:           p.HolesPlayed,
			Hole:            p.HolesPlayed,
			Today:           p.RoundScoreToPar,
			Total:           p.ScoreToPar,
			TotalStrokes:    totalStrokes,
			Rounds:          rounds,
		})
	}

	return &golf.Leaderboard{
		Tour:      euro.String(),
		TourIndex: euro.Index(),
		Players:   players,
	}, nil
}

func (euro *Euro) SetLeaderboard(lb *golf.Leaderboard) {
	euro.leaderboard = lb
}

func (euro *Euro) Leaderboard() *golf.Leaderboard {
	return euro.leaderboard
}

func (euro *Euro) LastUpdated() time.Time {
	return euro.lastUpdated
}

func (euro *Euro) SetLastUpdated(t time.Time) {
	euro.lastUpdated = t
}

func (euro *Euro) String() string {
	return "European Tour"
}

func (euro *Euro) Twitter() string {
	return "europeantour"
}

func (euro *Euro) TwitterID() string {
	return "55246492"
}

type EuroLeaderboard struct {
	Players []struct {
		PositionDesc    string
		ScoreToPar      int
		RoundScoreToPar int
		HolesPlayed     int
		FirstName       string
		LastName        string
		Rounds          []struct {
			Strokes int
		}
	}
}

func FixEuroName(s string) string {

	parts := strings.Split(s, ",")

	// unexpected format, just return the original string
	if len(parts) != 2 {
		return s
	}

	surname, name := parts[0], parts[1]

	var surnameParts []string

	for _, n := range strings.Split(surname, " ") {

		if n == "I" || n == "II" || n == "III" || n == "IV" {
			surnameParts = append(surnameParts, n)
			continue
		}

		fixed := strings.TrimSpace(string(n[0]) + strings.ToLower(n[1:]))
		surnameParts = append(surnameParts, fixed)
	}

	return strings.TrimSpace(name) + " " + strings.Join(surnameParts, " ")
}
