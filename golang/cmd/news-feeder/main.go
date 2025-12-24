package main

import (
	"fmt"
	"time"

	"github.com/heussd/nats-news-analysis/internal/langdetect"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/fulltextrss"
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

			retrievalStart := time.Now()
			fulltext := fulltextrss.RetrieveFullText(url)
			retrievalTime := time.Since(retrievalStart)

			news := model.MakeNews(fulltext)
			langPostfix := langdetect.AssignSubjectPostfixBasedOnLanguage(
				fmt.Sprintf("%s %s", fulltext.Title, fulltext.Content),
			)
			news.Language = string(langPostfix)

			if _, err := queue.Publish(
				news,
				queue.PublishSubject(
					fmt.Sprintf("news.%s", langPostfix)),
			); err != nil {
				logger.Error().Err(err).Msg("Failed to publish news")
			}

			logger.Info().
				Str("domain", fulltext.Domain).
				Int("fulltext-length", len(fulltext.Content)).
				Int64("retrieval-duration-ms", retrievalTime.Milliseconds()).
				Msg("Full text retrieval complete")
		},
		queue.SubscribeSubject("article-urls"),
		queue.StreamNameIsSubjectName(),
		queue.WaitTillSomeoneWants("news.undetermined"),
		queue.WaitTillSomeoneWants("news.de"),
		queue.WaitTillSomeoneWants("news.en"),
	); err != nil {
		panic(err)
	}
}
