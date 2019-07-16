package ss

import (
	"fmt"
	"io"
	"net/http"
)

func GetMatches(c http.Client) (io.ReadCloser, error) {

	res, err := c.Get("https://supersport.com/golf/results")
	if err != nil {
		return nil, err
	}
	return res.Body, nil

}

func FetchMatches(c http.Client) ([]Match, error) {

	r, err := GetMatches(c)
	if err != nil {
		return nil, err
	}
	return ParseMatches(r)
}

func GetMatch(c http.Client, id string) (io.ReadCloser, error) {

	res, err := c.Get(fmt.Sprintf("https://supersport.com/golf/match/%s", id))
	if err != nil {
		return nil, err
	}
	return res.Body, nil

}

func FetchMatch(c http.Client, matchId string) (Match, error) {

	r, err := GetMatch(c, matchId)
	if err != nil {
		return Match{}, err
	}
	return ParseMatch(r, matchId)
}

func GetScorecard(c http.Client, scorecardId string, matchId string) (io.ReadCloser, error) {

	url := fmt.Sprintf("https://www.supersport.com/api/golf/scorecard?fixture_id=%s&scorecard_id=%s", matchId, scorecardId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", "https://www.supersport.com/golf/")

	res, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	return res.Body, nil

}

func FetchScorecard(c http.Client, scorecardId string, matchId string) (Scorecard, error) {

	r, err := GetScorecard(c, scorecardId, matchId)
	if err != nil {
		return Scorecard{}, nil
	}
	return ParseScorecard(r)
}
