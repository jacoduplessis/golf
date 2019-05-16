package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://sunshinetour.com/?page_id=25983&tourn=LMPC&season=218S&report=http://sunshinetour.info/tic/tmscores.cgi?tourn=LMPC~season=218S

// https://sunshinetour.com/api/sst/cache/sst/218S/218S-LMPC-scores-latest.json?randomadd=1552215230156

type Sunshine struct {
	lastUpdated time.Time
	leaderboard *Leaderboard
	tid         string
}

func (ss *Sunshine) TID() string {
	return ss.tid
}

func (ss *Sunshine) Index() int {
	return 3
}

func (ss *Sunshine) UpdateTID() error {

	res, err := client.Get("https://sunshinetour.com/api/sst/cache/sst/tmticx")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var d struct {
		Code     string `json:"code"`
		TMParams struct {
			SeasonCode string `json:"season_code"`
		} `json:"tm_params"`
	}

	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return err
	}

	if d.Code == "" {
		return errors.New("Sunshine TID is empty")
	}
	ss.tid = d.Code + "," + d.TMParams.SeasonCode
	return nil
}

func (ss *Sunshine) Request() (*http.Request, error) {

	args := strings.Split(ss.TID(), ",")
	if len(args) < 2 {
		return nil, fmt.Errorf("[ss] Invalid TID %s", ss.TID())

	}
	tid, season := args[0], args[1]

	u := fmt.Sprintf("https://sunshinetour.com/api/sst/cache/sst/%s/%s-%s-scores-latest.json", season, season, tid)
	return http.NewRequest("GET", u, nil)
}

func (ss *Sunshine) Parse(r io.Reader) (*Leaderboard, error) {

	var d SunshineLeaderboard
	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	var players []*Player

	for _, p := range d.Result.Entry {

		var rounds []int

		for _, rs := range []string{p.R1, p.R2, p.R3, p.R4, p.R5, p.R6} {
			rounds = appendRound(rounds, rs)
		}

		var totalStrokes int
		for _, rd := range rounds {
			if rd > 40 {
				totalStrokes += rd
			}
		}

		total, _ := strconv.Atoi(p.Par)

		strokes, _ := strconv.Atoi(p.Score)

		hole, _ := strconv.Atoi(p.Hole)

		players = append(players, &Player{
			Name:            p.Name,
			Country:         p.Country,
			CurrentPosition: p.Position,
			Total:           total,
			Rounds:          rounds,
			TotalStrokes:    strokes,
			Hole:            hole,
			After:           hole,
		})
	}

	return &Leaderboard{
		Tour:       ss.String(),
		TourIndex:  ss.Index(),
		Tournament: d.Name,
		Course:     d.CourseName,
		Location:   fmt.Sprintf("%s, %s", d.CourseCity, d.CourseCountry),
		Players:    players,
	}, nil
}

func (ss *Sunshine) SetLeaderboard(lb *Leaderboard) {
	ss.leaderboard = lb
}

func (ss *Sunshine) Leaderboard() *Leaderboard {
	return ss.leaderboard
}

func (ss *Sunshine) LastUpdated() time.Time {
	return ss.lastUpdated
}

func (ss *Sunshine) SetLastUpdated(t time.Time) {
	ss.lastUpdated = t
}

func (ss *Sunshine) String() string {
	return "Sunshine Tour"
}

func (ss *Sunshine) Twitter() string {
	return "Sunshine_Tour"
}

func (ss *Sunshine) TwitterID() string {
	return "126255586"
}

type SunshineLeaderboard struct {
	Name          string `json:"short_name"`
	CourseName    string `json:"course_name"`
	CourseCity    string `json:"course_city"`
	CourseCountry string `json:"course_country"`

	Result struct {
		Entry []struct {
			Position string `json:"pos"`
			Score    string `json:"score"`
			Par      string `json:"vspar"`
			Name     string `json:"name"`
			Country  string `json:"nationality"`
			Hole     string `json:"holes"`
			R1       string `json:"score_R1"`
			R2       string `json:"score_R2"`
			R3       string `json:"score_R3"`
			R4       string `json:"score_R4"`
			R5       string `json:"score_R5"`
			R6       string `json:"score_R6"`
		} `json:"scores_entry"`
	} `json:"scores"`
}
