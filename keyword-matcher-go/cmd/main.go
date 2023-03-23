package main

import (
	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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
	var hostname, _ = os.LookupEnv("host")
	var logger = log.With().
		Str("service", "keyword-matcher-go").
		Str("container", hostname).
		Logger()

	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)
		var startTime = time.Now()
		var fulltext = fulltextrss.RetrieveFullText(url)
		var text = prepareAndCleanString(fulltext)

		var match, regexId = keywords.Match(text)
		var elapsedTime = time.Since(startTime)
		if match {
			queue.PushToPocket(model.Match{
				Url:     url,
				RegexId: regexId,
			})
		}

		logger.Info().
			Bool("match", match).
			Str("regex-id", regexId).
			Str("url", url).
			Int64("keyword-matching-duration-ms", elapsedTime.Milliseconds()).
			Msg("Analysis complete")

	})
	logger.Info().Msg("🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
