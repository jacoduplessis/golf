package golf

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
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

func GetListenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8000"
}

func AppendRound(rounds []int, rs string) []int {
	r, err := strconv.Atoi(rs)
	if err != nil {
		return rounds
	}
	if r > 40 {
		rounds = append(rounds, r)
	}
	return rounds
}
