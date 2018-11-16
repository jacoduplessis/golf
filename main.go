package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/jacoduplessis/twitterparse"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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
  <title>Golf Leaderboards</title>
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
	<p>View <a href="/news">news</a>.</p>
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
	newsTemplate = template.Must(template.New("item").Parse(`
			<div class="tweet">
				<p><strong>{{ .UserName }} (@{{ .UserHandle}})</strong> &middot; <time datetime="{{ .ISOTime }}">{{ .RelativeTime }}</time></p>
				
				<p>{{ .Content }}</p>
				{{ if .ImageURL }}
				<img src="{{ .ImageURL }}">
				{{end}}
			</div>
	`))

	// language=HTML format=true
	newsTemplate.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Golf News</title>
		<style>
			img {
				max-width: 100%;
			}

			.tweet {
				margin-top: 1rem;
				border: 2px solid #ccc;
				padding: 1rem;
				background-color: #fff;
			}
		</style>
	</head>
	<body style="max-width: 650px;margin: 0 auto; background-color: #eee">		
		
		<div style="text-align: center; margin-bottom: 3rem">
			<h1>Golf News</h1>
			<p>from Twitter</p>
		</div>
		
		
		{{ range . }}{{template "item" . }}{{end}}
		
		<p><a href="/">Home</a></p>
	</body>
</html>
`)
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

func news(w http.ResponseWriter, r *http.Request) *AppError {

	results := make(chan []*twitterparse.Tweet, len(tours))
	var tweets []*twitterparse.Tweet

	for _, tour := range tours {

		go func(t Tour) {
			tweets, err := twitterparse.FetchUserWithClientAndParse(client, t.Twitter())
			if err != nil {
				log.Printf("error fetching tweets for %s %s", t, err)
			}
			results <- tweets
		}(tour)

	}

	for range tours {
		tweets = append(tweets, <-results...)
	}

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].Timestamp > tweets[j].Timestamp
	})

	newsTemplate.ExecuteTemplate(w, "", tweets)
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

	noInitial := flag.Bool("no_initial", false, "No initial fetching of data")
	noInterval := flag.Bool("no_interval", false, "No periodic updates of data")

	flag.Parse()

	// initial data
	if !*noInitial {
		log.Println("fetching initial data")
		updateTournaments()
		updateLeaderboards()
	}

	// setup tickers
	if !*noInterval {
		log.Println("creating tickers")
		intervalUpdate()
	}

	parseTemplate()
	http.Handle("/", Handler(index))
	http.Handle("/news", Handler(news))
	addr := getListenAddr()
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
