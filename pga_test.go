package main

import (
	"bytes"
	"os"
	"testing"
)

func TestPGA(t *testing.T) {

	parseTemplate()
	pga := &PGA{}
	b, _ := os.Open("test/pga.json")
	out := &bytes.Buffer{}
	lb, err := pga.Parse(b)
	if err != nil {
		t.Fail()
	}
	err = tmpl.ExecuteTemplate(out, "leaderboard", lb)
	if err != nil {
		t.Fail()
	}
}
