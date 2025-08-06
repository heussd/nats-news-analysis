package keywords

import (
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/model"
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
