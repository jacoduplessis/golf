package main

import (
	"io/ioutil"
	"testing"
	"encoding/json"
	"os"
	"fmt"
)

func TestTemplate(t *testing.T) {

	b, _ := ioutil.ReadFile("leaderboard.json")

	var tournament Tournament
	err := json.Unmarshal(b, &tournament)
	fmt.Println(err)
	fmt.Println(tournament)
	tmpl.ExecuteTemplate(os.Stdout, "", &tournament)

}
