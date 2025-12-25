package htmlsanitise

import (
	"strings"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/microcosm-cc/bluemonday"
)

var bm = bluemonday.StrictPolicy()

func PrepareAndCleanString(news *model.News) string {
	text := strings.Join([]string{
		news.Title,
		news.Excerpt,
		bm.Sanitize(news.Content),
	}, " ")

	return text
}
