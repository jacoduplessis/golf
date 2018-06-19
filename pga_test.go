package main

import (
	"fmt"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {

	parseTemplate()
	pga := &PGA{}
	b, _ := os.Open("test_pga.json")
	out, _ := os.Create("test_pga.html")
	lb, _ := pga.Parse(b)
	err := tmpl.ExecuteTemplate(out, "leaderboard", lb)
	fmt.Println(err)
}
