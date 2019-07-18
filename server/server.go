package server

import (
	"github.com/jacoduplessis/golf"
	"github.com/jacoduplessis/golf/euro"
	"github.com/jacoduplessis/golf/pga"
	"github.com/jacoduplessis/golf/ss/server"
	"github.com/jacoduplessis/golf/sunshine"
	"github.com/jacoduplessis/twitterparse"
	"html/template"
	"net/http"
)

var tmpl *template.Template

var twitterClient = &twitterparse.TwitterClient{}
var newsTemplate *template.Template
var videoTemplate *template.Template

var tours = []golf.Tour{
	&pga.PGA{},
	&euro.Euro{},
	&sunshine.Sunshine{},
}

func GetHandler(c http.Client) http.Handler {

	parseTemplate()
	setupTwitterClient(c)
	updateTournaments(c, tours)
	updateLeaderboards(c, tours)
	intervalUpdate(c, tours)

	mux := http.NewServeMux()

	resultsHandler := server.GetHandler(c)

	mux.Handle("/", Handler(index))
	mux.Handle("/news", Handler(news))
	// mux.Handle("/videos", Handler(videosHandler))
	mux.Handle("/results/", http.StripPrefix("/results", resultsHandler))

	return mux
}

func GetServer(c http.Client) *http.Server {

	addr := golf.GetListenAddr()
	h := GetHandler(c)

	return &http.Server{
		Handler: h,
		Addr:    addr,
	}

}
