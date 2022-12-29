package config

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

var (
	UrlsFile               = GetEnv("URLS_FILE", "urls.txt")
	NatsServer             = GetEnv("NATS_SERVER", nats.DefaultURL)
	NatsOutputQueueName    = GetEnv("NATS_OUTPUT_QUEUE_NAME", "article-urls")
	NatsOutputQueueSubject = GetEnv("NATS_OUTPUT_QUEUE_SUBJECT", "article-url")
	NatsKeyValueBucket     = GetEnv("NATS_KV_BUCKET", "article-urls-proposed")
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
