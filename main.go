package main

import (
	"bytes"
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
	After           int // holes completed
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
		log.Println(apperr.Error)
	}
}

func renderTemplate(w http.ResponseWriter, data interface{}) *AppError {
	b := bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(&b, "", data); err != nil {
		return &AppError{"Unavailable, please try again in a minute.", 500, err}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	b.WriteTo(w)
	return nil
}

func renderJSON(w http.ResponseWriter, data interface{}) *AppError {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return &AppError{"Unavailable, please try again in a minute.", 500, err}
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
	<div style="margin-right: 1rem">	
		<h3>{{ .Tour }} - {{ .Tournament }}</h1>
		<h3>{{ .Course }}{{if .Location}}, {{ .Location }}{{end}}</h3>
		<table>
			<thead>
				<tr>
					<th title="position">Pos</th>
					<th>Par</th>
					<th>Player</th>
					<th title="nationality">Nat</th>
					<th title="current hole">On</th>
					<th title="holes played">Pl</th>
					<th title="this round">Rd</th>
					<th>Rounds</th>
					<th title="number of strokes">Total</th>
				</tr>
			</thead>
			<tbody>
			{{range .Players}}
				<tr>
					<td>{{ .CurrentPosition }}</td>
					<td style="text-align: right">{{ .Total }}</td>
					<td>{{ .Name }}</td>
					<td style="text-align: right">{{ .Country }}</td>
					<td style="text-align: right">{{ .Hole }}</td>
					<td style="text-align: right">{{ .After }}</td>
					<td style="text-align: right">{{ .Today }}</td>
					<td style="padding-left: 1rem">{{range .Rounds}}{{if .}}{{ . }} {{end}}{{end}}</td>
					<td style="text-align: right">{{ .TotalStrokes}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>
	</div>
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
	td:nth-of-type(3) { padding-left: 10px }
</style>
</head>
<body>
	<div style="display: flex; flex-flow: row wrap">
	{{range .Leaderboards}}{{template "leaderboard" .}}{{end}}
	</div>
	<p><a href="/?format=json">Get this data as JSON.</a></p>
</body>
</html>
`)
}

func updateTournaments() {

	results := make(chan *Leaderboard, len(tours))

	for _, tour := range tours {

		go func(t Tour) {
			r, err := t.Request()
			if err != nil {
				log.Printf("Error getting update request: %s\n", err)
				results <- nil
				return
			}
			resp, err := client.Do(r)
			if err != nil {
				log.Printf("Error performing update request: %s\n", err)
				results <- nil
				return
			}
			defer resp.Body.Close()
			lb, err := t.Parse(resp.Body)
			if err != nil {
				log.Printf("Error parsing leaderboard response: %s\n", err)
				results <- nil
				return
			}
			results <- lb

		}(tour)

	}

	for i := 0; i < len(tours); i++ {
		lb := <-results
		tour := tours[i]
		if lb != nil {
			tour.SetLeaderboard(lb)
			tour.SetLastUpdated(time.Now())
		}
	}

}

type TemplateContext struct {
	Leaderboards []*Leaderboard
}

func index(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := &TemplateContext{}

	for _, tour := range tours {
		ctx.Leaderboards = append(ctx.Leaderboards, tour.Leaderboard())
	}

	if r.URL.Query().Get("format") == "json" {
		return renderJSON(w, &ctx)
	}
	return renderTemplate(w, &ctx)

}

func intervalUpdate() {

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			updateTournaments()
		}
	}()
}

func main() {

	updateTournaments()
	intervalUpdate()
	parseTemplate()

	http.Handle("/", Handler(index))
	addr := getListenAddr()
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
