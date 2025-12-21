package main

import (
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/internal/raindrop"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var subject = utils.GetEnv("NATS_SUBJECT", "match-urls")

func main() {
	if err := queue.Subscribe(
		func(match *model.Match) {
			fmt.Printf("Received match from queue %s\n", match.Url)

			if err := raindrop.Add(match); err != nil {
				fmt.Printf("received error from raindrop: %w", err)
			} else {
				fmt.Printf("added to Raindrop: %s\n	", match.Url)
			}
		},
		queue.SubscribeSubject(subject),
		queue.StreamNameIsSubjectName(),
	); err != nil {
		panic(err)
	}
}
