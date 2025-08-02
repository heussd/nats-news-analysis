package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn
var js nats.JetStreamContext

func init() {

	for nc == nil {
		var err error
		nc, err = nats.Connect(config.NatsServer)
		if err != nil {
			fmt.Printf("Could not connect to NATS server at %s, because: ", config.NatsServer)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}
	}

	var err error
	js, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		panic(err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:       config.NatsOutputQueueName,
		Subjects:   []string{config.NatsOutputQueueSubject},
		Retention:  nats.LimitsPolicy,
		MaxAge:     time.Hour * 24 * 90,
		Duplicates: time.Hour * 24 * 30,
	})
	if err != nil {
		panic(err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:       config.NatsInputQueueName,
		Subjects:   []string{config.NatsInputQueueSubject},
		Retention:  nats.LimitsPolicy,
		MaxAge:     time.Hour * 24 * 90,
		Duplicates: time.Hour * 24 * 30,
	})
	if err != nil {
		panic(err)
	}

	var consumer *nats.ConsumerInfo
	consumer, err = js.AddConsumer(
		config.NatsInputQueueName,
		&nats.ConsumerConfig{
			Durable:       config.NewsStreamConsumer,
			DeliverPolicy: nats.DeliverNewPolicy,
		})
	if err != nil {
		fmt.Println("There was a problem creating the consumer. This can occur if the consumer already exists but with a different config.")
		panic(err)
	} else {
		println("Using consumer:", consumer.Name)
	}

}

func WithArticleUrls(f func(m *nats.Msg)) {
	sub, err := js.PullSubscribe(
		config.NatsInputQueueSubject,
		config.NewsStreamConsumer,
	)
	if err != nil {
		fmt.Printf("Cannot subscribe to subject %s with consumer %s, because: ",
			config.NatsInputQueueName,
			config.NewsStreamConsumer)
		fmt.Println(err)
	}

	for {
		msgs, err := sub.Fetch(
			5,
			nats.MaxWait(60*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			log.Printf("Fetch error: %v", err)
			panic(err)
		}

		for _, msg := range msgs {
			msg.Ack()
			f(msg)
		}
	}
}

func PushToPocket(match model.Match) {
	headers := nats.Header{}
	headers.Add(nats.MsgIdHdr, "news-keyword-matcher"+match.Url)

	data, _ := json.Marshal(match)

	msg := &nats.Msg{
		Header:  headers,
		Subject: config.NatsOutputQueueSubject,
		Data:    data,
	}

	js.PublishMsg(msg)
}
