package main

import (
	"time"

	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/fulltextrss"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		input    = queue.AddStreamOrDie(utils.GetEnv("NATS_INPUT_STREAM", "article-urls"), queue.DefaultDupeWindow)
		output   = queue.AddStreamOrDie(utils.GetEnv("NATS_OUTPUT_STREAM", "news"), queue.DefaultDupeWindow)
		consumer = queue.AddConsumerOrDie(input, utils.GetEnv("NATS_CONSUMER", "default"))
	)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().
		Str("service", "news-feeder-go").
		Logger()

	err := queue.Subscribe(input, consumer,
		func(article *model.Article) {
			url := article.Url

			retrievalStart := time.Now()
			fulltext := fulltextrss.RetrieveFullText(url)
			retrievalTime := time.Since(retrievalStart)

			queue.Publish(output,
				model.News{
					Title:    fulltext.Title,
					Excerpt:  fulltext.Excerpt,
					Author:   fulltext.Author,
					Language: fulltext.Language,
					URL:      fulltext.Url,
					Content:  fulltext.Content,
					Date:     fulltext.Date,
				},
				"", false,
			)

			logger.Info().
				Str("domain", fulltext.Domain).
				Int("fulltext-length", len(fulltext.Content)).
				Int64("retrieval-duration-ms", retrievalTime.Milliseconds()).
				Msg("Full text retrieval complete")
		}, true)
	if err != nil {
		panic(err)
	}
}
