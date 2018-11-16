package main

import (
	"bytes"
	"os"
	"testing"
)

func TestEuro(t *testing.T) {

	parseTemplate()
	euro := &Euro{}
	b, _ := os.Open("test/euro.json")
	out := &bytes.Buffer{}
	lb, err := euro.Parse(b)
	if err != nil {
		t.Fail()
	}
	err = tmpl.ExecuteTemplate(out, "leaderboard", lb)
	if err != nil {
		t.Fail()
	}
}
