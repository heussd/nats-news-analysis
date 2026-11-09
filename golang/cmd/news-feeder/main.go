package main

import (
	"fmt"
	"time"

	"github.com/heussd/nats-news-analysis/internal/langdetect"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/defuddle"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "news-feeder-go").
		Logger()

	if err := queue.Subscribe(
		func(article *model.Article) {
			url := article.Url
			var news model.News
			var err error

			retrievalStart := time.Now()
			if news, err = defuddle.RetrieveFullText(url); err != nil {
				logger.Error().Err(err).Str("url", url).Msg("Failed to retrieve full text for url")
				return
			}
			retrievalTime := time.Since(retrievalStart)

			langPostfix := langdetect.AssignSubjectPostfixBasedOnLanguage(
				news.Content,
			)
			news.Language = string(langPostfix)

			if _, err := queue.Publish(
				news,
				queue.PublishSubject(fmt.Sprintf("news.%s", langPostfix)),
				queue.MessageID(url),
			); err != nil {
				logger.Error().Err(err).Msg("Failed to publish news")
			}

			logger.Info().
				Str("domain", news.Domain).
				Int("fulltext-length", len(news.Content)).
				Int64("retrieval-duration-ms", retrievalTime.Milliseconds()).
				Msg("News feeding complete")
		},
	); err != nil {
		panic(err)
	}
}
