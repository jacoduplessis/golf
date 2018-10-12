package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Sunshine struct {
	lastUpdated time.Time
	leaderboard *Leaderboard
	tid         string
}

func (ss *Sunshine) TID() string {
	return ss.tid
}

func (ss *Sunshine) UpdateTID() error {

	res, err := client.Get("https://sunshinetour.com/api/sst/cache/sst/tmticx")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var d struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return err
	}

	if d.Code == "" {
		return errors.New("Sunshine TID is empty")
	}
	ss.tid = d.Code
	return nil
}

func (ss *Sunshine) Request() (*http.Request, error) {

	u := fmt.Sprintf("http://sunshinetour.com/api/sst/cache/sst/218S/218S-%s-result-result-PF.json", ss.TID())
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

		players = append(players, &Player{
			Name:            p.Name,
			Country:         p.Country,
			CurrentPosition: p.Position,
			Total:           total,
			Rounds:          rounds,
			TotalStrokes:    strokes,
		})
	}

	return &Leaderboard{
		Tour:       ss.String(),
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
	return "sunshinetour"
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
			R1       string `json:"score_R1"`
			R2       string `json:"score_R2"`
			R3       string `json:"score_R3"`
			R4       string `json:"score_R4"`
			R5       string `json:"score_R5"`
			R6       string `json:"score_R6"`
		} `json:"result_entry"`
	}
}
