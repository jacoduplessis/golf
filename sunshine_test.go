package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSunshineParser(t *testing.T) {

	parseTemplate()
	ss := &Sunshine{}
	b, _ := os.Open("test_sunshine.json")
	out, _ := os.Create("test_sunshine.html")
	lb, _ := ss.Parse(b)
	err := tmpl.ExecuteTemplate(out, "leaderboard", lb)
	fmt.Println(err)
}
