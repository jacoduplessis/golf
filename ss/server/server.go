package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jacoduplessis/golf/ss"
	"net/http"
	"os"
	"time"
)

var cacheControl = "public, max-age=86400" // one day

func renderTemplateWithHeaders(w http.ResponseWriter, name string, data interface{}, headers map[string]string) error {
	tmpl, ok := tmpl[name]
	if !ok {
		return fmt.Errorf("template %s does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encountered.
	buf := bytes.Buffer{}

	err := tmpl.ExecuteTemplate(&buf, "", data)
	if err != nil {
		fmt.Printf("template error in '%s': %v\n", name, err)
		return err
	}

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	_, err = buf.WriteTo(w)
	return err
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	headers := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": cacheControl,
	}
	return renderTemplateWithHeaders(w, name, data, headers)
}

func renderError(w http.ResponseWriter, code int, message string, err error) {
	fmt.Printf("error: %s %v", message, err)
	w.WriteHeader(code)
	_, _ = w.Write([]byte(message))
}

type APIHandler struct {
	client http.Client
}

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	matches, err := ss.FetchMatches(h.client)
	if err != nil {
		renderError(w, 500, "Upstream Data Error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", cacheControl)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	_ = json.NewEncoder(w).Encode(matches)

}

type ResultsHandler struct {
	client http.Client
}

func (h ResultsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	matches, err := ss.FetchMatches(h.client)
	if err != nil {
		renderError(w, 500, "Upstream Data Error", err)
		return
	}

	_ = renderTemplate(w, "index", matches)

}

type TournamentHandler struct {
	client http.Client
}

func (h TournamentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	matchId := vars["matchId"]
	match, err := ss.FetchMatch(h.client, matchId)

	context := []ss.Match{
		match,
	}

	if err != nil {
		renderError(w, 500, "Upstream Data Error", err)
		return
	}

	_ = renderTemplate(w, "index", context)
}

type ScorecardHandler struct {
	client http.Client
}

func (h ScorecardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	matchId := vars["matchId"]
	scorecardId := vars["scorecardId"]

	if matchId == "" || scorecardId == "" {
		renderError(w, 400, "Bad Request", nil)
	}

	sc, err := ss.FetchScorecard(h.client, scorecardId, matchId)
	if err != nil {
		renderError(w, 500, "Upstream Data Error", err)
		return
	}

	_ = renderTemplate(w, "scorecard", sc)

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

func GetHandler(c http.Client) http.Handler {
	r := mux.NewRouter()

	r.Handle("/", ResultsHandler{client: c})
	r.Handle("/api", APIHandler{client: c})
	r.Handle("/tournaments/{matchId}", TournamentHandler{client: c})
	r.Handle("/scorecards/{matchId}/{scorecardId}", ScorecardHandler{client: c})

	return r
}

func GetServer(c http.Client) *http.Server {

	h := GetHandler(c)

	return &http.Server{
		Handler:      h,
		Addr:         getListenAddr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
