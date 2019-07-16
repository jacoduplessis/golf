package ss

import (
	"os"
	"testing"
)

func TestParseMatches(t *testing.T) {

	r, _ := os.Open("../test/ss_results.html")

	matches, err := ParseMatches(r)
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 9 {
		t.Error("Expected 4 matches from results page")
	}
}

func TestParseMatch(t *testing.T) {
	r, _ := os.Open("../test/ss_match.html")

	match, err := ParseMatch(r, "123")

	if err != nil {
		t.Fatal(err)
	}

	if len(match.Players) != 70 {
		t.Error("Expected 70 players")
	}

	first := match.Players[0]

	if first.Name != "Dylan Frittelli" {
		t.Error("Wrong name")
	}

	if first.MatchId != "123" {
		t.Error("MatchID injected value not present")
	}

}

func TestParseScorecard(t *testing.T) {

	r, _ := os.Open("../test/ss_scorecard.json")

	sc, err := ParseScorecard(r)

	if err != nil {
		t.Fatal(err)
	}

	if len(sc.Rounds) != 4 {
		t.Error("Expected 4 rounds")
	}
}
