package nats

import (
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

var nc *nats.Conn
var js nats.JetStreamContext

func init() {
	var err error
	nc, err = nats.Connect(config.NatsServer)
	if err != nil {
		panic(err)
	}

	js, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		panic(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:       config.NatsOutputQueueName,
		Subjects:   []string{config.NatsOutputQueueSubject},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	})

}

func WithArticleUrls(f func(m *nats.Msg)) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err := js.QueueSubscribe(config.NatsInputQueueSubject, config.NatsInputQueueName, f)
	if err != nil {
		panic(err)
	}

	wg.Done()
	wg.Wait()
}

func PushToPocket(s string) {
	headers := nats.Header{}
	headers.Add(nats.MsgIdHdr, "news-keyword-matcher"+s)

	msg := &nats.Msg{
		Header:  headers,
		Subject: config.NatsOutputQueueSubject,
		Data:    []byte(s),
	}

	js.PublishMsg(msg)
}
