package main

import (
	"fmt"
	"html/template"
	"strings"
)

func URLize(s string) template.HTML {

	words := strings.Split(s, " ")
	var processed []string

	for _, word := range words {
		if strings.HasPrefix(word, "http://") {
			wordText := word
			if len(word) > 20 {
				wordText = fmt.Sprintf("%s...", word[:20])
			}
			processed = append(processed, fmt.Sprintf(`<a href="%s" target="_blank" referrerpolicy="no-referrer">%s</a>`, word, wordText))
			continue
		}
		processed = append(processed, word)
	}

	return template.HTML(strings.Join(processed, " "))
}

func FixEuroName(s string) string {

	parts := strings.Split(s, ",")

	// unexpected format, just return the original string
	if len(parts) != 2 {
		return s
	}

	surname, name := parts[0], parts[1]

	var surnameParts []string

	for _, n := range strings.Split(surname, " ") {

		if n == "I" || n == "II" || n == "III" || n == "IV" {
			surnameParts = append(surnameParts, n)
			continue
		}

		fixed := strings.TrimSpace(string(n[0]) + strings.ToLower(n[1:]))
		surnameParts = append(surnameParts, fixed)
	}

	return strings.TrimSpace(name) + " " + strings.Join(surnameParts, " ")
}
