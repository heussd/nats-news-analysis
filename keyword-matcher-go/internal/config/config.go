package config

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

var (
	KeywordsFileUrl        = GetEnv("KEYWORDS_FILE_URL", "https://raw.githubusercontent.com/heussd/nats-news-analysis/refs/heads/main/keyword-matcher-go/internal/keywords/keywords.txt")
	NewsStreamConsumer     = GetEnv("NEWS_STREAM_CONSUMER", "default")
	FullTextRssServer      = GetEnv("FULLTEXTRSS_SERVER", "http://localhost:80")
	NatsServer             = GetEnv("NATS_SERVER", nats.DefaultURL)
	NatsInputQueueName     = GetEnv("NATS_INPUT_QUEUE_NAME", "article-urls")
	NatsInputQueueSubject  = GetEnv("NATS_INPUT_QUEUE_SUBJECT", "article-url")
	NatsOutputQueueName    = GetEnv("NATS_OUTPUT_QUEUE_NAME", "match-urls")
	NatsOutputQueueSubject = GetEnv("NATS_OUTPUT_QUEUE_SUBJECT", "match-url")
)

// GetEnv Taken from https://stackoverflow.com//questions/40326540/how-to-assign-default-value-if-env-var-is-empty#answer-45978733
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	fmt.Println(key, "=", value)
	return value
}
