package keywords

import (
	"fmt"
	"os"

	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/cloudtextfile"
)

func Match(s string) (bool, string) {
	var keywords []model.KeywordEntry
	var err error

	if keywords, err = cloudtextfile.CachedParsedKeywords(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range keywords {
		if match, _ := v.Regexp.MatchString(s); match {
			return true, v.Id
		}
	}
	return false, ""
}
