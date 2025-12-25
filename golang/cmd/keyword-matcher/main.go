package main

import (
	"time"

	"github.com/heussd/nats-news-analysis/internal/htmlsanitise"
	"github.com/heussd/nats-news-analysis/internal/keywords"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	subscribeSubject = utils.GetEnv("NATS_SUBSCRIBE_SUBJECT", "news.*")
	consumerName     = utils.GetEnv("NATS_CONSUMER", "default")
	publishSubject   = utils.GetEnv("NATS_PUBLISH_SUBJECT", "match-urls")
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "keyword-matcher-go").
		Logger()

	err := queue.Subscribe(
		func(news *model.News) {
			text := htmlsanitise.PrepareAndCleanString(news)

			matchingStart := time.Now()
			var matched []model.Keyword
			var err error
			if matched, err = keywords.Match(text); err != nil {
				logger.Error().Err(err).Msg("Failed to match keywords")
				return
			}
			matchingTime := time.Since(matchingStart)

			if len(matched) > 0 {
				_, err := queue.Publish(
					model.Match{
						Url:      news.URL,
						Keywords: matched,
					},
					queue.PublishSubject(publishSubject),
				)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to publish match")
				}
			}

			logger.Info().
				Bool("match", len(matched) > 0).
				Int("fulltext-length", len(text)).
				Int64("keyword-matching-duration-ms", matchingTime.Milliseconds()).
				Msg("Analysis complete")
		},
		queue.SubscribeSubject(subscribeSubject),
		queue.SubscribeConsumer(consumerName),
	)
	if err != nil {
		panic(err)
	}
}
