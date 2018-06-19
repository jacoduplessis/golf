package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AppError struct {
	Message string
	Code    int
	Error   error
}

type Player struct {
	Name            string
	Country         string
	CurrentPosition string
	StartPosition   string
	Through         int // holes completed
	Hole            int // current hole
	Today           int // par
	Total           int // par
	TotalStrokes    int // stroke
	Rounds          []int
}

type Leaderboard struct {
	Tour       string
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
	SetLeaderboard(*Leaderboard)
	Leaderboard() *Leaderboard
	LastUpdated() time.Time
	SetLastUpdated(time.Time)
}

var tours = []Tour{
	&PGA{},
	&Euro{},
}

var tmpl *template.Template

var client = http.Client{Timeout: time.Second * 10}

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

func parseTemplate() {
	// language=HTML format=true
	tmpl = template.Must(template.New("leaderboard").Parse(`
		<h1>{{ .Tour }} - {{ .Tournament }}</h1>
		{{if .Location}}
		<h3>{{ .Location }}</h3>
		{{end}}
		<h3>{{ .Course }}</h3>
		{{if .Date}}
		<h4>{{ .Date }}</h4>
		{{end}}
		<p>Current Round: {{ .Round }}{{if .Updated}} | Updated: {{ .Updated }}{{end}}</p>
		
		<table>
			<thead>
				<tr>
					<th>Position</th>
					<th>Started</th>
					<th>Player</th>
					<th>Country</th>
					<th>Hole</th>
					<th>Through</th>
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
					<td>{{ .StartPosition }}</td>
					<td>{{ .Name }}</td>
					<td style="text-align: right">{{ .Country }}</td>
					<td style="text-align: right">{{ .Hole }}</td>
					<td style="text-align: right">{{ .Through }}</td>
					<td style="text-align: right">{{ .Total }}</td>
					<td style="text-align: right">{{ .Today }}</td>
					<td style="padding-left: 1rem">{{range .Rounds}}{{ . }} {{end}}</td>
					<td style="text-align: right">{{ .TotalStrokes}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>
	`))

	// language=HTML format=true
	tmpl.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>Golf</title>
  <style>
	html { font-family: monospace }
	table { border-collapse: collapse; border: 1px solid }
	th { text-align: left; padding-left: 0.4rem; padding-right: 0.4rem }
	tr:hover { background-color: #ccc }
</style>
</head>
<body>
	{{range .Leaderboards}}{{template "leaderboard" .}}{{end}}
</body>
</html>
`)
}

func updateTournaments() {

	errs := make(chan error, len(tours))

	for _, tour := range tours {

		if time.Now().Sub(tour.LastUpdated()) < time.Second*60 {
			errs <- nil
			continue
		}

		fmt.Println("Updating", tour)

		go func(t Tour) {
			r, err := t.Request()
			if err != nil {
				errs <- err
				return
			}
			resp, err := client.Do(r)
			if err != nil {
				errs <- err
				return
			}
			defer resp.Body.Close()
			lb, err := t.Parse(resp.Body)
			if err != nil {
				errs <- err
				return
			}
			t.SetLeaderboard(lb)
			t.SetLastUpdated(time.Now())
			errs <- nil
		}(tour)

	}

	for i := 0; i < len(tours); i++ {
		e := <-errs
		if e != nil {
			fmt.Println(e)
		}
	}

}

type TemplateContext struct {
	Leaderboards []*Leaderboard
}

func index(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := &TemplateContext{}

	updateTournaments()
	for _, tour := range tours {
		ctx.Leaderboards = append(ctx.Leaderboards, tour.Leaderboard())
	}

	if r.URL.Query().Get("format") == "json" {
		return renderJSON(w, &ctx)
	}
	return renderTemplate(w, &ctx)

}

func main() {

	parseTemplate()
	http.Handle("/", HandlerFunc(index))
	addr := getListenAddr()
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
