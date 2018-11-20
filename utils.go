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
