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
	Leaderboard() *Leaderboard
	SetLeaderboard(*Leaderboard)
	LastUpdated() time.Time
	SetLastUpdated(time.Time)
	TID() string
	UpdateTID() error
	Twitter() string
}

var tours = []Tour{
	&PGA{},
	&Euro{},
	&Sunshine{},
}

var tmpl *template.Template

var client = http.Client{Timeout: time.Second * 10}

var newsTemplate *template.Template

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
		<h3>{{ .Tour }} - {{ .Tournament }}</h3>
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
	<p>Click on any country code to highlight all players from that country.</p>
	<p><a href="/?format=json">Get this data as JSON.</a></p>

	<script>
		(function(){
			const cells = document.querySelectorAll('td:nth-of-type(4)') // country code
			cells.forEach(el => {
				el.addEventListener('click', event => {					
					cells.forEach(e => e.parentElement.style.backgroundColor = '') // reset all
					cells.forEach(e => {
						if (e.textContent === event.target.textContent) e.parentElement.style.backgroundColor = 'yellow' // highlight all from same country
					})	
				})
			})
		})()
	</script>
</body>
</html>
`)

	// language=HTML
	newsTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html></html>
	`))
}

func updateTournaments() {

	errs := make(chan error, len(tours))

	for _, tour := range tours {
		go func(t Tour) {
			errs <- t.UpdateTID()
		}(tour)
	}

	for range tours {
		err := <-errs
		if err != nil {
			log.Printf("Error updating TID: %s", err)
		}
	}
}

func updateLeaderboards() {

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

	for _, tour := range tours {
		lb := <-results
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

type NewsFeed struct {
	TourName string
	Items    []struct {
		Title   string
		Content string
	}
}

func news(w http.ResponseWriter, r *http.Request) *AppError {

	results := make(chan *NewsFeed, len(tours))

	for _, tour := range tours {

		go func(t Tour) {
			resp, err := client.Get(fmt.Sprintf("http://www.twitrss.me/twitter_user_to_rss/?user=%s", tour.Twitter()))
			if err != nil {
				results <- nil
				return
			}
			defer resp.Body.Close()

			// parse response into NewsFeed
			results <- &NewsFeed{}
		}(tour)

	}

	var feeds []NewsFeed
	for i := 0; i < len(tours); i++ {
		feed := <-results
		feeds = append(feeds, feed)
	}

	newsTemplate.ExecuteTemplate(w, "", feeds)
	return nil
}

func intervalUpdate() {

	tickerLeaderboards := time.NewTicker(1 * time.Minute)
	go func() {
		for range tickerLeaderboards.C {
			updateLeaderboards()
		}
	}()

	tickerTournaments := time.NewTicker(6 * time.Hour)
	go func() {
		for range tickerTournaments.C {
			updateTournaments()
		}
	}()

}

func main() {

	// initial data
	updateTournaments()
	updateLeaderboards()

	// setup tickers
	intervalUpdate()

	parseTemplate()
	http.Handle("/", Handler(index))
	addr := getListenAddr()
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
