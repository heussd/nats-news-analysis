package main

import (
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/internal/raindrop"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

func main() {
	var (
		stream   = queue.AddStreamOrDie(utils.GetEnv("NATS_INPUT_STREAM", "match-urls"), queue.DefaultDupeWindow)
		consumer = queue.AddConsumerOrDie(stream, utils.GetEnv("NATS_CONSUMER", "default"))
	)

	var err = queue.Subscribe(stream, consumer,
		func(match *model.Match) {
			fmt.Printf("Received match from queue %s\n", match.Url)

			if err := raindrop.Add(match); err != nil {
				fmt.Printf("received error from raindrop: %w", err)
			} else {
				fmt.Printf("added to Raindrop: %s\n	", match.Url)
			}
		},
		true,
	)

	if err != nil {
		panic(err)
	}
}
