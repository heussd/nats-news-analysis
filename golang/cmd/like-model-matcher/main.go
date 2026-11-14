package main

import (
	"fmt"
	"time"

	"github.com/heussd/nats-news-analysis/internal/likemodel"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	subscribeSubject = utils.GetEnv("NATS_SUBSCRIBE_SUBJECT", "news.*")
	consumerName     = utils.GetEnv("NATS_CONSUMER", "like-model-matcher")
	publishSubject   = utils.GetEnv("NATS_PUBLISH_SUBJECT", "match-urls")
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "like-model-matcher-go").
		Logger()

	err := queue.Subscribe(
		func(news *model.News) {
			text := fmt.Sprintf("%s %s %s", news.Title, news.Description, news.ContentRawText)

			matchingStart := time.Now()
			var liked bool
			var err error
			if liked, err = likemodel.Classify(text); err != nil {
				logger.Error().Err(err).Msg("Failed to classify text")
				return
			}
			matchingTime := time.Since(matchingStart)

			if liked {
				_, err := queue.Publish(
					model.Match{
						Url: news.URL,
						Keywords: []model.Keyword{
							{
								Text: "Like Model",
							},
						},
					},
					queue.PublishSubject(publishSubject),
					queue.MessageID(news.URL),
				)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to publish match")
				}
			}

			logger.Info().
				Bool("match", liked).
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
