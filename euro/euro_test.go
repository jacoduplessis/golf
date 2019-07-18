package euro

import "testing"

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
