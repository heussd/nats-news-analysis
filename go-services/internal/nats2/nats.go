package nats

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	"nats-news-analysis/internal/config"
)

var js nats.JetStreamContext

func init() {

	for nc == nil {
		if nc, err = nats.Connect(config.NatsServer); err != nil {
			fmt.Printf("Could not connect to NATS server at %s, because: ", config.NatsServer)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}
	}

	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
		fmt.Println("HERE")
		panic(err)
	}

	if _, err = js.AddStream(&nats.StreamConfig{
		Name: config.NatsOutputQueueName,
		Subjects: []string{
			config.NatsOutputQueueSubject,
		},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	}); err != nil {

		panic(err)
	}

}

func WithFeedUrls(f func(m *nats.Msg)) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var subscribed = false
	for !subscribed {
		_, err := js.QueueSubscribe(config.NatsInputQueueSubjectFeed, config.NatsInputQueueNameFeed, f)
		if err != nil {
			fmt.Printf("Cannot subscribe queue %s subject %s, because: ", config.NatsInputQueueNameFeed, config.NatsInputQueueSubjectFeed)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			subscribed = true
		}
	}

	wg.Done()
	wg.Wait()
}
