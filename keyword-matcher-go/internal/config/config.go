package config

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

var (
	KeywordsFile           = GetEnv("KEYWORDS_FILE", "keywords.txt")
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
