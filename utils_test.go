package main

import "testing"

func TestURLize(t *testing.T) {

	txt := "Link to http://example.com"

	html := URLize(txt)

	if html != `Link to <a href="http://example.com" target="_blank" referrerpolicy="no-referrer">http://example.com</a>` {
		t.Errorf("Incorrect markup returned: %s", html)
	}

	txt = "A very long http://1234567890123456790.link.com"

	html = URLize(txt)

	if html != `A very long <a href="http://1234567890123456790.link.com" target="_blank" referrerpolicy="no-referrer">http://1234567890123...</a>` {
		t.Errorf("Incorrect markup returned: %s", html)
	}

}

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
