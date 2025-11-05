package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/heussd/nats-news-analysis/internal/keywords"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bm = bluemonday.StrictPolicy()

func prepareAndCleanString(news *model.News) string {
	text := strings.Join([]string{
		news.Title,
		news.Excerpt,
		bm.Sanitize(news.Content),
	}, " ")

	return text
}

func main() {
	var (
		input    = queue.AddStreamOrDie(utils.GetEnv("NATS_INPUT_STREAM", "news"), queue.DefaultDupeWindow)
		output   = queue.AddStreamOrDie(utils.GetEnv("NATS_OUTPUT_STREAM", "match-urls"), queue.DefaultDupeWindow)
		consumer = queue.AddConsumerOrDie(input, utils.GetEnv("NATS_CONSUMER", "default"))
	)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "keyword-matcher-go").
		Logger()

	err := queue.Subscribe(input, consumer,
		func(news *model.News) {
			text := prepareAndCleanString(news)

			matchingStart := time.Now()
			var matched []model.Keyword
			var err error
			if matched, err = keywords.Match(text); err != nil {
				fmt.Printf("error matching: %w", err)
			}
			matchingTime := time.Since(matchingStart)

			if len(matched) > 0 {
				queue.Publish(output,
					model.Match{
						Url:      news.URL,
						Keywords: matched,
					},
					"", false,
				)
			}

			logger.Info().
				Bool("match", len(matched) > 0).
				Int("fulltext-length", len(text)).
				Int64("keyword-matching-duration-ms", matchingTime.Milliseconds()).
				Msg("Analysis complete")
		}, true)
	if err != nil {
		panic(err)
	}
}
