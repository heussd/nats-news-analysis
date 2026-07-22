package main

import (
	"strings"
	"time"

	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/internal/ngrams"
	"github.com/heussd/nats-news-analysis/internal/timeseries"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	subscribeSubject = utils.GetEnv("NATS_SUBSCRIBE_SUBJECT", "news.*")
	consumerName     = "ngram-analysis"
	minimumNGramSize = 1
	maximumNGramSize = 4
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "ngram-analysis").
		Logger()

	err := queue.Subscribe(
		func(news *model.News) {
			if !strings.HasPrefix(news.Language, "en") {
				return
			}

			ngrams, err := ngrams.ParseAndGenerateStatistics(news, minimumNGramSize, maximumNGramSize)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate n-gram statistics")
				return
			}

			for i := range ngrams {
				ngrams[i].Source = news.URL
				ngrams[i].Language = news.Language
				ngrams[i].Timestamp = news.Date

				if ngrams[i].Timestamp == "" || timeseries.ValidateTimestamp(ngrams[i].Timestamp) != nil {
					// Fallback directly in the n-gram field when timestamp is missing or invalid.
					ngrams[i].Timestamp = time.Now().UTC().Format("2006-01-02 15:04:05-07")
				}
			}

			if err := timeseries.AddTimeSeriesData(ngrams); err != nil {
				logger.Error().Err(err).Msg("Failed to add n-gram statistics to time series database")
				return
			}

			logger.Info().
				Msg("N-gram statistics generated")
		},
		queue.SubscribeSubject(subscribeSubject),
		queue.SubscribeConsumer(consumerName),
	)
	if err != nil {
		panic(err)
	}
}
