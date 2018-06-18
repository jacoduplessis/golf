package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"html/template"
	"fmt"
	"time"
	"errors"
)

type AppError struct {
	Message string
	Code    int
	Error   error
}

// language=HTML format=true
var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>Golf</title>
  <style>
	table { border-collapse: collapse; border: 1px solid;}
	th { text-align: left; padding-left: 0.4rem; padding-right: 0.4rem;}
	tr:hover { background-color: #ccc; }
</style>
</head>
<body>
	<p>Updated: {{.LastUpdated}}</p>
	{{with .Leaderboard }}
		<h2>{{ .TourName }} - {{ .TournamentName }}</h2>
		<h3>{{range .Courses}}{{ .CourseName }} {{end}}</h3>
		<h4>{{ .StartDate }} - {{ .EndDate }}</h4>
	
		<p>Current Round: {{ .CurrentRound }}</p>
		
		<table>
			<thead>
				<tr>
					<th>Position</th>
					<th>Player</th>
					<th>Country</th>
					<th>Thru</th>
					<th>Total</th>
					<th>Today</th>
					<th>Rounds</th>
					<th>Strokes</th>
				</tr>
			</thead>
			<tbody>
			{{range .Players}}
				<tr>
					<td>{{ .CurrentPosition }}</td>
					<td>{{ .PlayerBio.FirstName }} {{ .PlayerBio.LastName }}</td>
					<td>{{ .PlayerBio.Country }}</td>
					<td>{{ .Thru }}</td>
					<td>{{ .Total }}</td>
					<td style="text-align: center">{{ .Today }}</td>
					<td>{{range .Rounds}}{{ .Strokes }} {{end}}</td>
					<td style="text-align: right">{{ .TotalStrokes}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>
	{{end}}
</body>
</html>
`))

var client = http.Client{Timeout: time.Second * 10}
var tournament Tournament
var lastUpdated time.Time

type Handler func(w http.ResponseWriter, r *http.Request) *AppError

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if apperr := fn(w, r); apperr != nil {
		http.Error(w, apperr.Message, apperr.Code)
		log.Println(apperr.Error.Error())
	}
}

func HandlerFunc(fn func(w http.ResponseWriter, r *http.Request) *AppError) Handler {
	return fn
}

func renderTemplate(w http.ResponseWriter, data interface{}) *AppError {

	if err := tmpl.ExecuteTemplate(w, "", data); err != nil {
		return &AppError{"Template error", 500, err}
	}
	return nil
}

func renderJSON(w http.ResponseWriter, data interface{}) *AppError {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return &AppError{"Could not write response", 500, err}
	}
	return nil
}

func getListenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8000"
}

type Tournament struct {
	LastUpdated string `json:"last_updated"`
	Leaderboard struct {
		Courses []struct {
			CourseName string `json:"course_name"`
		}
		TournamentName string `json:"tournament_name"`
		TourName       string `json:"tour_name"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`

		CurrentRound int `json:"current_round"`

		Players []struct {
			CourseHole      int    `json:"course_hole"`
			CurrentPosition string `json:"current_position"`
			StartPosition   string `json:"start_position"`
			Thru            int
			Today           int
			Total           int
			TotalStrokes    int    `json:"total_strokes"`
			PlayerBio struct {
				Country   string `json:"country"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				ShortName string `json:"short_name"`
			} `json:"player_bio"`
			Rounds []struct {
				RoundNumber int `json:"round_number"`
				Strokes     int `json:"strokes"`
			}
		} `json:"players"`
	} `json:"leaderboard"`
}

func updateTournament() error {
	if time.Now().Sub(lastUpdated) < time.Second * 60 {
		return nil // still fresh
	}
	var current struct {
		TID string `json:"tid"`
	}
	resp, err := client.Get("https://statdata.pgatour.com/r/current/message.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&current); err != nil {
		return err
	}

	if current.TID == "" {
		return errors.New("TID is empty")
	}
	resp, err = client.Get(fmt.Sprintf("https://statdata.pgatour.com/r/%s/leaderboard-v2mini.json", current.TID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	lastUpdated = time.Now()
	return json.NewDecoder(resp.Body).Decode(&tournament)

}

func index(w http.ResponseWriter, r *http.Request) *AppError {

	err := updateTournament()
	if err != nil {
		return &AppError{Code: 500, Message: "Error parsing upstream data", Error: err}
	}

	if r.URL.Query().Get("format") == "json" {
		return renderJSON(w, &tournament)
	}
	return renderTemplate(w, &tournament)

}

func main() {

	http.Handle("/", HandlerFunc(index))
	addr := getListenAddr()
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
