package main

import (
	"strings"
	"time"

	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bm = bluemonday.StrictPolicy()

func prepareAndCleanString(fulltext fulltextrss.RSSFullTextResponse) string {
	var text = strings.Join([]string{
		fulltext.Title,
		fulltext.Excerpt,
		bm.Sanitize(fulltext.Content),
	}, " ")

	return text
}

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var logger = log.With().
		Str("service", "keyword-matcher-go").
		Logger()

	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)

		var retrievalStart = time.Now()
		var fulltext = fulltextrss.RetrieveFullText(url)
		var retrievalTime = time.Since(retrievalStart)

		var text = prepareAndCleanString(fulltext)

		var matchingStart = time.Now()
		var match, regexId = keywords.Match(text)
		var matchingTime = time.Since(matchingStart)

		if match {
			queue.PushToPocket(model.Match{
				Url:     url,
				RegexId: regexId,
			})
		}

		logger.Info().
			Bool("match", match).
			Str("regex-id", regexId).
			Str("domain", fulltext.Domain).
			Int("fulltext-length", len(text)).
			Int64("retrieval-duration-ms", retrievalTime.Milliseconds()).
			Int64("keyword-matching-duration-ms", matchingTime.Milliseconds()).
			Msg("Analysis complete")

	})

}
