package htmlsanitise

import (
	"html"
	"strings"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/microcosm-cc/bluemonday"
)

func ExtractFields(news *model.News) string {
	text := strings.Join([]string{
		news.Title,
		news.Excerpt,
		news.Content,
	}, " ")

	return text
}

func Sanitize(text string) string {
	text = bluemonday.
		StrictPolicy().
		AddSpaceWhenStrippingTag(true).
		Sanitize(text)

	text = html.UnescapeString(text)
	text = strings.Join(strings.Fields(text), " ")
	text = strings.TrimSpace(text)

	return text
}
