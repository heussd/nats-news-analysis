package cloudtextfile

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
)

var (
	cachedKeywords []model.KeywordEntry
	lastGenerated  time.Time
	cacheDuration  = 20 * time.Minute
)

func RetrieveKeywordsFile() (keywords []string, err error) {
	client := &http.Client{}
	var req *http.Request

	if req, err = http.NewRequest("GET", config.KeywordsFileUrl, nil); err != nil {
		return nil, err
	}

	var response *http.Response

	if response, err = client.Do(req); err != nil {
		return nil, err
	}

	if status := response.StatusCode; status != 200 {
		return nil, fmt.Errorf("failed to retrieve keywords file: status code %d", status)
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	keywords = strings.Split(string(body), "\n")

	return keywords, nil
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

var cleanUpRegexes = []regexp2.Regexp{
	*regexp2.MustCompile("[^a-zA-Z]", 0),
	*regexp2.MustCompile("\\s\\S\\s", 0),
	*regexp2.MustCompile("\\s\\s+", 0),
}

func CachedParsedKeywords() (keywords []model.KeywordEntry, err error) {
	if time.Since(lastGenerated) > cacheDuration {

		keywords = []model.KeywordEntry{}

		var plainKeywords []string
		if plainKeywords, err = RetrieveKeywordsFile(); err != nil {
			return nil, err
		}

		for _, text := range plainKeywords {
			if text == "" || strings.HasPrefix(text, "#") {
				continue
			}

			fmt.Printf("Parsing \"%s\" as regex\n", text)

			var regex = regexp2.MustCompile(text, 0)
			keywords = append(keywords, model.KeywordEntry{
				Regexp: *regex,
				Id:     humanReadable(text),
				Text:   text,
			})
		}
		cachedKeywords = keywords
		lastGenerated = time.Now()
	}

	return cachedKeywords, nil
}
