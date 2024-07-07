package config

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

var (
	FullTextRssServer             = GetEnv("FULLTEXTRSS_SERVER", "http://localhost:80")
	KeywordsFile                  = GetEnv("KEYWORDS_FILE", "keywords.txt")
	NatsInputQueueNameArticle     = GetEnv("NATS_INPUT_QUEUE_NAME_ARTICLE", "article-urls")
	NatsInputQueueNameFeed        = GetEnv("NATS_INPUT_QUEUE_NAME_FEED", "feed-urls")
	NatsInputQueueNameMatch       = GetEnv("NATS_INPUT_QUEUE_NAME_MATCH", "match-urls")
	NatsInputQueueSubjectArticle  = GetEnv("NATS_INPUT_QUEUE_SUBJECT_ARTICLE", "article-url")
	NatsInputQueueSubjectFeed     = GetEnv("NATS_INPUT_QUEUE_NAME_FEED", "feed-url")
	NatsInputQueueSubjectMatch    = GetEnv("NATS_INPUT_QUEUE_NAME_MATCH", "match-url")
	NatsKeyValueBucket            = GetEnv("NATS_KV_BUCKET", "article-urls-proposed")
	NatsOutputQueueNameArticle    = GetEnv("NATS_OUTPUT_QUEUE_NAME_ARTICLE", "article-urls")
	NatsOutputQueueNameMatch      = GetEnv("NATS_OUTPUT_QUEUE_NAME_MATCH", "match-urls")
	NatsOutputQueueSubjectArticle = GetEnv("NATS_OUTPUT_QUEUE_SUBJECT_ARTICLE", "article-url")
	NatsOutputQueueSubjectMatch   = GetEnv("NATS_OUTPUT_QUEUE_SUBJECT_MATCH", "match-url")
	NatsServer                    = GetEnv("NATS_SERVER", nats.DefaultURL)
	RaindropCollection            = GetEnv("RAINDROP_COLLECTION", "")
	RaindropToken                 = GetEnv("RAINDROP_TOKEN", "")
	UrlsFile                      = GetEnv("URLS_FILE", "urls.txt")
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
