package keywords

import (
	"bufio"
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"os"
	"regexp"
	"strings"
)

var keywords []regexp.Regexp = parseKeywordsFile()

func parseKeywordsFile() []regexp.Regexp {
	var keywords []regexp.Regexp

	readFile, err := os.Open(config.KeywordsFile)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		var text = fileScanner.Text()

		if text == "" ||
			strings.HasPrefix(text, "#") {
			continue
		}

		fmt.Printf("Parsing \"%s\" as regex\n", text)

		var regex = regexp.MustCompile(text)
		keywords = append(keywords, *regex)
	}

	return keywords
}

func Match(s string) bool {
	for i, v := range keywords {
		if i%25 == 0 {
			fmt.Printf(" ... keywords >= %d\n", i)
		}
		if v.MatchString(s) {
			fmt.Println(" ✅ Match")
			return true
		}
	}
	fmt.Println(" ❌ No match")
	return false
}
