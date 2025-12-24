package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/heussd/nats-news-analysis/internal/cloudtextfile"
	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	urls     = utils.GetEnv("URLS_URL", "")
	waitTime = time.Hour * 2
)

func main() {
	for {
		lines, _, err := cloudtextfile.CachedCloudTextFile(urls)
		if err != nil {
			panic(err)
		}

		for _, line := range lines {

			stripped := strings.TrimSpace(line)
			if stripped == "" {
				continue
			}

			u, err := url.ParseRequestURI(stripped)
			if err != nil || u.Scheme == "" || u.Host == "" {
				println("invalid url:", stripped)
				continue
			}

			fmt.Printf("Publishing url %s\n", u.String())
			if _, err = nats.Publish(
				model.Feed{
					Url: u.String(),
				},
				func(npo *nats.NatsPublishOptions) {
					npo.Subject = "feed-urls"
					npo.NatsMessageID = u.String()
				},
			); err != nil {
				panic(err)
			}
		}

		fmt.Printf("Waiting for %s\n", waitTime)
		time.Sleep(waitTime)
	}
}
