package keywords

import (
	"bufio"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"os"
	"strings"
)

type KeywordEntry struct {
	regexp regexp2.Regexp
	id     string
	text   string
}

var cleanUpRegexes = []regexp2.Regexp{
	*regexp2.MustCompile("[^a-zA-Z]", 0),
	*regexp2.MustCompile("\\s\\S\\s", 0),
	*regexp2.MustCompile("\\s\\s+", 0),
}

var keywords []KeywordEntry

func init() {
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

		var regex = regexp2.MustCompile(text, 0)
		keywords = append(keywords, KeywordEntry{
			regexp: *regex,
			id:     humanReadable(text),
			text:   text,
		})
	}

	if len(keywords) == 0 {
		fmt.Println("Error: No keywords found")
		os.Exit(1)
	}
}

func humanReadable(regex string) string {
	var s = regex
	var err error
	for _, r := range cleanUpRegexes {
		if s, err = r.Replace(s, " ", 0, -1); err != nil {
			panic(err)
		}
	}

	return strings.TrimSpace(s)
}

func Match(s string) (bool, string) {
	for _, v := range keywords {
		if match, _ := v.regexp.MatchString(s); match {
			return true, v.id
		}
	}
	return false, ""
}
