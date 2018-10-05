package main

import (
	"fmt"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"os"
	"testing"
)

func TestIDExtraction(t *testing.T) {

	var tid string
	markup, _ := os.Open("euro_home.html")
	defer markup.Close()
	root, _ := html.Parse(markup)
	s := cascadia.MustCompile("#ETContainer_thisWeek>div")
	d := s.MatchFirst(root)
	for _, a := range d.Attr {
		if a.Key == "id" {
			tid = a.Val
		}
	}
	fmt.Println(tid)
}
