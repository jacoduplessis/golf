package euro

import (
	"fmt"
	"os"
	"testing"
)

func TestFixEuroName(t *testing.T) {

	s, expected := "CABRERA BELLO, Rafa", "Rafa Cabrera Bello"
	r := FixEuroName(s)
	if r != expected {
		t.Errorf("Expected %s, got %s", expected, r)
	}

	s, expected = "HOWELL III, Charles", "Charles Howell III"
	r = FixEuroName(s)
	if r != expected {
		t.Errorf("Expected %s, got %s", expected, r)
	}

}

func TestTIDExtraction(t *testing.T) {

	r, _ := os.Open("../test/euro.html")
	defer r.Close()
	tid, _ := extractTID(r)
	expected := "2019090"

	if tid != expected {
		t.Errorf("Expected TID to be %s, got %s", expected, tid)
	}
}

func TestParse(t *testing.T) {
	r, _ := os.Open("../test/euro.json")
	e := &Euro{}
	lb, err := e.Parse(r)
	fmt.Printf("%s %#v", err, lb)
}
