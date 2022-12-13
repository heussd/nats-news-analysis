package main

import (
	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/nats-io/nats.go"
	"strings"
)

func main() {
	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)
		var fulltext = fulltextrss.RetrieveFullText(url)

		var text = strings.Join([]string{
			fulltext.Title,
			fulltext.Excerpt,
			fulltext.Content,
		}, " ")

		if keywords.Match(text) {
			queue.PushToPocket(url)
		}
	})
}
