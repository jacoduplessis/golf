package euro

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/cascadia"
	"github.com/jacoduplessis/golf"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// http://www.europeantour.com/library/MyEt/TourType=1/leaderboardJSON.js
// http://www.europeantour.com/data/tournament/2018050/leaderboard/languagecode/eng/isproam/0/feed/

// can also get tid here
// http://app.europeantour.com/mobile/

type Euro struct {
	lastUpdated time.Time
	leaderboard *golf.Leaderboard
	tid         string
}

var tidSelector = cascadia.MustCompile("#ETContainer_thisWeek>div")

func (euro *Euro) TID() string {
	return euro.tid
}

func (euro *Euro) Index() int {
	return 2
}

func (euro *Euro) UpdateTID(c http.Client) error {

	res, err := c.Get("http://www.europeantour.com")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var tid string
	root, err := html.Parse(res.Body)
	if err != nil {
		return err
	}

	d := tidSelector.MatchFirst(root)
	for _, a := range d.Attr {
		if a.Key == "id" {
			tid = a.Val
		}
	}
	if tid == "" {
		return errors.New("Euro TID is empty")
	}
	euro.tid = tid
	return nil
}

func (euro *Euro) Request() (*http.Request, error) {

	u := fmt.Sprintf("http://www.europeantour.com/data/tournament/%s/leaderboard/languagecode/eng/isproam/0/feed/", euro.TID())
	return http.NewRequest("GET", u, nil)
}

func (euro *Euro) Parse(r io.Reader) (*golf.Leaderboard, error) {

	var d EuroLeaderboard
	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	var players []*golf.Player
	var tournamentName string

	for i, p := range d.LeaderboardData {
		if i == 0 {
			tournamentName = p.CupSeasonName
		}
		today, _ := strconv.Atoi(p.Today)
		total, _ := strconv.Atoi(p.Topar)
		after, _ := strconv.Atoi(p.HolesPlayed)
		hole, _ := strconv.Atoi(p.Hole)
		var rounds []int

		for _, rs := range []string{p.R1, p.R2, p.R3, p.R4, p.R5, p.R6} {
			rounds = golf.AppendRound(rounds, rs)
		}

		var totalStrokes int
		for _, rd := range rounds {
			if rd > 40 {
				totalStrokes += rd
			}
		}

		players = append(players, &golf.Player{
			Name:            FixEuroName(p.Name),
			Country:         p.Countrycode,
			Today:           today,
			CurrentPosition: p.Position,
			Total:           total,
			After:           after,
			Hole:            hole,
			Rounds:          rounds,
			TotalStrokes:    totalStrokes,
		})
	}
	meta := d.LeaderboardCourseInfoData[0]
	return &golf.Leaderboard{
		Tour:       euro.String(),
		TourIndex:  euro.Index(),
		Tournament: tournamentName,
		Course:     meta.CourseName,
		Location:   fmt.Sprintf("%s, %s", meta.CityName, meta.CountryName),
		Players:    players,
	}, nil
}

func (euro *Euro) SetLeaderboard(lb *golf.Leaderboard) {
	euro.leaderboard = lb
}

func (euro *Euro) Leaderboard() *golf.Leaderboard {
	return euro.leaderboard
}

func (euro *Euro) LastUpdated() time.Time {
	return euro.lastUpdated
}

func (euro *Euro) SetLastUpdated(t time.Time) {
	euro.lastUpdated = t
}

func (euro *Euro) String() string {
	return "European Tour"
}

func (euro *Euro) Twitter() string {
	return "europeantour"
}

func (euro *Euro) TwitterID() string {
	return "55246492"
}

type EuroLeaderboard struct {
	LeaderboardCourseInfoData []struct {
		CityName    string // "Southampton, New York",
		CountryName string // "USA",
		CourseName  string // "Shinnecock Hills GC",
	}
	LeaderboardData []struct {
		// "LeaderboardData": [
		//    {
		//      "Amount": "0",
		//      "CountryName": "USA",
		Countrycode string //      "Countrycode": "USA",
		//      "CourseInstanceColor1": "",
		//      "CourseInstanceColor2": "",
		//      "CourseInstanceColor3": "",
		//      "CourseInstanceColor4": "",
		//      "CourseInstanceColor5": "",
		//      "CourseInstanceColor6": "",
		//      "CourseInstanceId1": "3652",
		//      "CourseInstanceId2": "3652",
		//      "CourseInstanceId3": "3652",
		//      "CourseInstanceId4": "3652",
		//      "CourseInstanceId5": "",
		//      "CourseInstanceId6": "",
		//      "CourseType1": "",
		//      "CourseType2": "",
		//      "CourseType3": "",
		//      "CourseType4": "",
		//      "CourseType5": "",
		//      "CourseType6": "",
		CupSeasonName string //      "CupSeasonName": "US OPEN",
		//      "CurrentRound": "4",
		//      "Current_Order": "1",
		//      "CutValue": "60",
		//      "Diff": "-",
		//      "EventStatusId": "104",
		//      "GroupId": "33",
		//      "H1": "4",
		//      "H10": "35",
		//      "H11": "33",
		//      "H12": "3",
		//      "H13": "4",
		//      "H14": "4",
		//      "H15": "4",
		//      "H16": "4",
		//      "H17": "4",
		//      "H18": "4",
		//      "H2": "2",
		//      "H3": "3",
		//      "H4": "4",
		//      "H5": "4",
		//      "H6": "5",
		//      "H7": "3",
		//      "H8": "4",
		//      "H9": "4",
		//      "HIn": "3",
		//      "HOut": "5",
		//      "HasBallSponsor": false,
		Hole string //      "Hole": "18",
		//      "HoleN": "18",
		HolesPlayed string //      "HolesPlayed": "18",
		//      "HolesPlayedMobile": "18",
		//      "IdCupKind": "1",
		//      "IsQualified": "False",
		//      "LastModifyDate": "6/18/2018 2:21:46 AM",
		Name string //      "Name": "KOEPKA, Brooks",
		//      "NoOfRounds": "4",
		//      "PlayerCurrentRound": "4",
		//      "PlayerId": "38783",
		//      "PlayerStatus": "P",
		//      "PointBasedTournament": "False",
		//      "Pos": "1",
		//      "PosMobile": "1",
		Position string //      "Position": "1",
		R1       string
		R2       string
		R3       string
		R4       string
		R5       string
		R6       string
		//      "R1": "75",
		//      "R1N": "75",
		//      "R1ToPar": "5       ",
		//      "R2": "66",
		//      "R2D": "-",
		//      "R2DN": "999",
		//      "R2N": "66",
		//      "R2ToPar": "-4      ",
		//      "R3": "72",
		//      "R3N": "72",
		//      "R3ToPar": "2       ",
		//      "R4": "68",
		//      "R4N": "68",
		//      "R4ToPar": "-2      ",
		//      "R5": "",
		//      "R5N": "",
		//      "R5ToPar": "",
		//      "R6": "",
		//      "R6N": "",
		//      "R6ToPar": "",
		//      "RankingPrizeMoney": "2,160,000.00",
		//      "ScoringPlayerStatusId": "0",
		//      "ScoringType": "Hole By Hole",
		//      "ScoringTypeTag": "lb_hole_by_hole",
		//      "Section": "1",
		//      "SectionOrder": "1",
		//      "StartN": "1",
		//      "StartTee": "1",
		//      "Starts": "T1",
		//      "TeeTime": "14:13",
		//      "TheOpenChampionship": "False",
		//      "Title": "2018",
		Today string //      "Today": "-2",
		//      "TodayMobile": "-2",
		//      "TodayN": "-2",
		Topar string //      "Topar": "1",
		//      "ToparL": "1",
		Total string //      "Total": "281",
		//      "WR": "4",
		//      "WRN": "4",
		//      "ispro": "True",
		//      "wCompCupSeasonID": "2018050"
	}
}

func FixEuroName(s string) string {

	parts := strings.Split(s, ",")

	// unexpected format, just return the original string
	if len(parts) != 2 {
		return s
	}

	surname, name := parts[0], parts[1]

	var surnameParts []string

	for _, n := range strings.Split(surname, " ") {

		if n == "I" || n == "II" || n == "III" || n == "IV" {
			surnameParts = append(surnameParts, n)
			continue
		}

		fixed := strings.TrimSpace(string(n[0]) + strings.ToLower(n[1:]))
		surnameParts = append(surnameParts, fixed)
	}

	return strings.TrimSpace(name) + " " + strings.Join(surnameParts, " ")
}
