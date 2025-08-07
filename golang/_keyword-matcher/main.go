package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/heussd/nats-news-analysis/internal/keywords"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/fulltextrss"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
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
	var (
		input    = queue.AddStreamOrDie(utils.GetEnv("NATS_INPUT_STREAM", "article-urls"), queue.DefaultDupeWindow)
		output   = queue.AddStreamOrDie(utils.GetEnv("NATS_OUTPUT_STREAM", "match-urls"), queue.DefaultDupeWindow)
		consumer = queue.AddConsumerOrDie(input, utils.GetEnv("NATS_CONSUMER", "default"))
	)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var logger = log.With().
		Str("service", "keyword-matcher-go").
		Logger()

	var err = queue.Subscribe(input, consumer,
		func(article *model.Article) {
			var url = article.Url

			var retrievalStart = time.Now()
			var fulltext = fulltextrss.RetrieveFullText(url)
			var retrievalTime = time.Since(retrievalStart)

			var text = prepareAndCleanString(fulltext)

			var matchingStart = time.Now()
			var matched []model.Keyword
			var err error
			if matched, err = keywords.Match(text); err != nil {
				fmt.Printf("error matching: %w", err)
			}
			var matchingTime = time.Since(matchingStart)

			if len(matched) > 0 {
				queue.Publish(output,
					model.Match{
						Url:      article.Url,
						Keywords: matched,
					},
					"", false,
				)
			}

			logger.Info().
				Bool("match", len(matched) > 0).
				Str("domain", fulltext.Domain).
				Int("fulltext-length", len(text)).
				Int64("retrieval-duration-ms", retrievalTime.Milliseconds()).
				Int64("keyword-matching-duration-ms", matchingTime.Milliseconds()).
				Msg("Analysis complete")

		}, true)

	if err != nil {
		panic(err)
	}
}
