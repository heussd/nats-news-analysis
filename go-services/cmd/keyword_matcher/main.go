package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"nats-news-analysis/internal/keywords"
	"nats-news-analysis/internal/model"
	queue "nats-news-analysis/internal/nats"
	"nats-news-analysis/pkg/fulltextrss"
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
	logger.Info().Msg("🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
