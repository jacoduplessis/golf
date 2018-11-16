package main

import (
	"bytes"
	"os"
	"testing"
)

func TestSunshine(t *testing.T) {

	parseTemplate()
	ss := &Sunshine{}
	b, _ := os.Open("test/sunshine.json")
	out := &bytes.Buffer{}
	lb, err := ss.Parse(b)
	if err != nil {
		t.Fail()
	}
	err = tmpl.ExecuteTemplate(out, "leaderboard", lb)
	if err != nil {
		t.Fail()
	}
}
