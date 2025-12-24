package keywords

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/heussd/nats-news-analysis/internal/cloudtextfile"
	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	cachedKeywords  []model.Keyword
	keywordsFileUrl = utils.GetEnv("KEYWORDS_FILE_URL", "https://raw.githubusercontent.com/heussd/nats-news-analysis/refs/heads/main/golang/internal/keywords/keywords.txt")
)

func Match(s string) (matched []model.Keyword, err error) {
	var keywords []model.Keyword
	if keywords, err = CachedParsedKeywords(); err != nil {
		return nil, fmt.Errorf("cannot get keywords: %w", err)
	}

	for _, v := range keywords {
		if match, _ := v.Regexp.MatchString(s); match {
			matched = append(matched, v)
		}
	}
	return matched, nil
}

func humanReadable(regex string) string {
	s := regex
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

func CachedParsedKeywords() (keywords []model.Keyword, err error) {
	plainKeywords, fromCache, err := cloudtextfile.CachedCloudTextFile(keywordsFileUrl)
	if err != nil {
		return nil, fmt.Errorf("cannot use cachedcloudtextfile %w", err)
	}

	if fromCache {
		return cachedKeywords, nil
	}

	keywords = []model.Keyword{}

	for _, text := range plainKeywords {
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}

		fmt.Printf("Parsing \"%s\" as regex\n", text)

		regex := regexp2.MustCompile(text, 0)
		keywords = append(keywords, model.Keyword{
			Regexp: *regex,
			Id:     humanReadable(text),
			Text:   text,
		})
	}
	cachedKeywords = keywords

	return cachedKeywords, nil
}
